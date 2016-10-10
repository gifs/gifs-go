package gifs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/odeke-em/semalim"
)

var (
	ErrExpectingAtLeastOneSource = errors.New("expecting atleast one source")
	ErrNilParamDereference       = errors.New("nil params dereference")

	errUnimplemented  = errors.New("unimplemented")
	errIllogicalState = errors.New("illogical and unexpected state")

	// If set, enables debug logging.
	debug = os.Getenv("DEBUG_GIFS_PKG") != ""
)

const (
	importEndpointURL             = "https://api.gifs.com/media/import"
	defaultConcurrentImportsCount = 10
)

var (
	debugLogPrintf func(string, ...interface{}) = func() func(string, ...interface{}) {
		if !debug {
			// noop
			return func(s string, args ...interface{}) {}
		}
		return func(s string, args ...interface{}) {
			log.Printf("[gifs] "+s, args...)
		}
	}()
)

type Client struct {
	client *http.Client
	apiKey string
}

type Option interface {
	apply(*Client)
}

type withAPIKey string

func (k withAPIKey) apply(g *Client) {
	g.apiKey = string(k)
}

func WithAPIKey(key string) Option {
	return withAPIKey(key)
}

type withClient struct {
	hc *http.Client
}

func (wc withClient) apply(g *Client) {
	g.client = wc.hc
}

func WithHTTPClient(hc *http.Client) Option {
	return withClient{hc}
}

func New(opts ...Option) (*Client, error) {
	c := &Client{}
	for _, o := range opts {
		o.apply(c)
	}
	return c, nil
}

type Trim struct {
	Start float64 `json:"start,omitempty"`
	End   float64 `json:"end,omitempty"`
}

type Request struct {
	// CreatedFrom when set tells the api whom
	// it should tag the caller as for purposes
	// of metric tracking and aggregation.
	CreatedFrom string `json:"caller,omitempty"`

	Title string `json:"title,omitempty"`

	// URL is the HTTP based URI pointing to
	// the media that is being transcoded.
	URL string `json:"source,omitempty"`

	// APIKey associates an authenticated user with
	// their imported media for purposes of ownership,
	// deals, special requests and many other good things.
	// If you don't have an API key, you can get one from
	// https://gifs.com/dashboard/api
	APIKey string `json:"api_key,omitempty"`

	// Tags are used for categorization of media
	Tags []string `json:"tags,omitempty"`

	// NSFW if set tells the API that this media is
	// Not-Safe-For-Work or Not-Suitable-For-Work
	NSFW bool `json:"nsfw,omitempty"`

	// Trim defines how to clip/trim the input source media
	// for example to trim the media from 10.45s to 22.3s
	// set the start to 10.45 and end to 22.3
	Trim *Trim `json:"trim,omitempty"`

	// Only set media if you are performing an upload
	media io.Reader

	// Attribution makes an association and gives
	// credit to the creator of the media.
	Attribution *Attribution `json:"attribution,omitempty"`

	// Crop defines a rectangle of area of interest that the output media
	// should be made first before any transcoding is done,
	// referenced by an (x, y) situated at the top left corner of the rectangle.
	// Note:
	// * The value of (x+width) must be less than or equal to the width of the media
	// * The value of (y+height) must be less than or equal to the height of the media
	Crop *Crop `json:"crop,omitempty"`

	callbackURI string `json:"-"`
}

func (p *Request) SetMedia(r io.Reader) error {
	if p == nil {
		return ErrNilParamDereference
	}
	p.media = r
	return nil
}

type FilesMap map[string]string

type responseError struct {
	Message string `json:"message,omitempty"`
}

func (re *responseError) Error() string {
	if re == nil {
		return ""
	}
	return re.Message
}

func (re *responseError) MarshalJSON() ([]byte, error) {
	if re == nil {
		return nil, nil
	}
	quoted := strconv.Quote(re.Message)
	return []byte(quoted), nil
}

func (re *responseError) UnmarshalJSON(bs []byte) error {
	re.Message = string(bs)
	return nil
}

type wrapperResponse struct {
	Success *Response      `json:"success,omitempty"`
	Errors  *responseError `json:"errors,omitempty"`
}

type Response struct {
	Embed  string        `json:"embed,omitempty"`
	OEmbed string        `json:"oembed,omitempty"`
	Files  FilesMap      `json:"files,omitempty"`
	Page   string        `json:"page,omitempty"`
	Error  responseError `json:"error,omitempty"`
}

func (res Response) HasFiles() bool {
	return len(res.Files) >= 1
}

func (res Response) File(mt MediaType) string {
	return res.Files[mt.Extension()]
}

func (g *Client) Upload() (*Response, error) {
	return nil, errUnimplemented
}

// Import is a method with which you'll specify atleast
// an http based URL pointing to media that you'd like
// to import to gifs.com.
func (g *Client) Import(req *Request) (*Response, error) {
	if req == nil {
		return nil, ErrNilParamDereference
	}

	bip := &BulkImportRequest{
		Requests: []*Request{req},
	}
	responses, err := g.ImportBulk(bip)
	if err != nil {
		return nil, err
	}
	if len(responses) < 1 {
		return nil, errIllogicalState
	}
	return responses[0], nil
}

// ImportSources is a convenience method that allows you to just specify
// multiple media URLs without having to construct each `Request` object.
func (g *Client) ImportSources(sources ...string) ([]*Response, error) {
	if len(sources) < 1 {
		return nil, ErrExpectingAtLeastOneSource
	}

	preparedRequest := []*Request{}
	for _, source := range sources {
		preparedRequest = append(preparedRequest, &Request{URL: source})
	}

	bip := &BulkImportRequest{Requests: preparedRequest}
	return g.ImportBulk(bip)
}

type BulkImportRequest struct {
	ConcurrentImports uint
	Requests          []*Request
}

func (p *Request) transformToImportBody() ([]byte, error) {
	if p == nil {
		return nil, ErrNilParamDereference
	}

	return json.Marshal(p)
}

func copyHeaders(from, to http.Header) {
	for key, _ := range from {
		fromValues := from[key]
		for _, value := range fromValues {
			to.Add(key, value)
		}
	}
}

func (g *Client) doPOSTRequest(uri string, req *Request, headers http.Header) (*http.Response, error) {
	b, err := req.transformToImportBody()
	if err != nil {
		return nil, err
	}
	debugLogPrintf("body %s for req: %+v", b, req)
	httpReq, err := http.NewRequest("POST", uri, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	copyHeaders(headers, httpReq.Header)
	httpReq.Header.Set("Content-Type", "application/json")
	if g.apiKey != "" {
		httpReq.Header.Set("Gifs-Api-Key", g.apiKey)
	}

	return g.httpClient().Do(httpReq)
}

func (g *Client) doMultipartUpload() (*http.Response, error) {
	return nil, errUnimplemented
}

func (g *Client) httpClient() *http.Client {
	if g.client != nil {
		return g.client
	}
	return http.DefaultClient
}

type jobType uint

const (
	postRequest jobType = iota
	uploadRequest
)

type httpRequestJob struct {
	uri     string
	req     *Request
	headers http.Header

	uuid uint64
	g    *Client
	typ  jobType
}

func (hj httpRequestJob) Id() interface{} {
	return hj.uuid
}

func (hj httpRequestJob) Do() (interface{}, error) {
	res, err := hj.g.doPOSTRequest(hj.uri, hj.req, hj.headers)
	debugLogPrintf("id: %v httpResposne: %v err: %v\n", hj.uuid, res, err)
	if err != nil {
		return nil, err
	}

	slurp, err := ioutil.ReadAll(res.Body)
	debugLogPrintf("id: %v slurp: %s err: %v\n", hj.uuid, slurp, err)
	if err != nil {
		return nil, err
	}
	_ = res.Body.Close()
	wrapperRes := new(wrapperResponse)
	err = json.Unmarshal(slurp, wrapperRes)
	debugLogPrintf("id: %v after unmarshalling, got %v err: %v\n", hj.uuid, wrapperRes, err)
	if err != nil {
		return nil, err
	}
	return wrapperRes, nil
}

func categorizeParallelJobResponses(resultsChan chan semalim.Result, maxResponseId uint64) ([]*Response, error) {
	idList := []uint64{}
	idMap := make(map[uint64]*Response)

	for result := range resultsChan {
		res, err, id := result.Value(), result.Err(), result.Id()
		debugLogPrintf("id: %d res: %v err: %v", id, res, err)

		var idKey uint64
		switch v := id.(type) {
		case int:
			idKey = uint64(v)
		case uint64:
			idKey = v
		case int64:
			idKey = uint64(v)
		default:
			parsedI, err := strconv.ParseUint(fmt.Sprintf("%s", v), 10, 64)
			if err != nil {
				// TODO: Log this to the user?
				// Otherwise we don't want to mess up our unique results
				// Shouldn't happen but if it does alas
				continue
			}
			idKey = parsedI
		}

		var finalRes *Response
		if res != nil {
			wrapRes := res.(*wrapperResponse)
			if wrapRes != nil {
				finalRes = wrapRes.Success
				if err == nil {
					err = wrapRes.Errors
				}
			}
		}

		if finalRes == nil {
			finalRes = new(Response)
			if err != nil {
				finalRes.Error = responseError{Message: err.Error()}
			}
		}
		idMap[idKey] = finalRes
		idList = append(idList, idKey)
	}

	debugLogPrintf("idMap: %v\n", idMap)
	// Now we've got to sort the results in the order that their requests were initially prepared
	responsesList := make([]*Response, maxResponseId)

	sort.Sort(uint64Slice(idList))
	for _, id := range idList {
		debugLogPrintf("\n\nid: %v v: %v\n\n", id, idMap[id])
		responsesList[id] = idMap[id]
	}

	return responsesList, nil
}

type uint64Slice []uint64

func (u64s uint64Slice) Len() int           { return len(u64s) }
func (u64s uint64Slice) Less(i, j int) bool { return u64s[i] < u64s[j] }
func (u64s uint64Slice) Swap(i, j int) {
	u64s[i], u64s[j] = u64s[j], u64s[i]
}

// ImportBulk is a convenience method that helps you import multiple media
// in one pass, however import requests will be made in parallel to
// the API. Responses per request will be matched by index/order of the requests.
func (g *Client) ImportBulk(bip *BulkImportRequest) ([]*Response, error) {
	var concurrentImports uint64 = defaultConcurrentImportsCount
	if bip.ConcurrentImports > 0 {
		concurrentImports = uint64(bip.ConcurrentImports)
	}

	maxResponseId := uint64(len(bip.Requests))
	jobsBench := make(chan semalim.Job)
	go func() {
		defer close(jobsBench)
		for i := uint64(0); i < maxResponseId; i++ {
			req := bip.Requests[i]
			jobsBench <- httpRequestJob{uri: importEndpointURL, req: req, uuid: uint64(i), g: g}
		}
	}()

	resultsChan := semalim.Run(jobsBench, concurrentImports)
	return categorizeParallelJobResponses(resultsChan, maxResponseId)
}
