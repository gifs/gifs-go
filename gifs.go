package gifs_go

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
	"sync"

	"github.com/odeke-em/semalim"
)

var (
	ErrExpectingAtLeastOneSource = errors.New("expecting atleast one source")
	ErrNilParamDereference       = errors.New("nil params dereference")

	errUnimplemented  = errors.New("unimplemented")
	errIllogicalState = errors.New("illogical and unexpected state")
)

type MediaType uint

const (
	MP4 MediaType = 1 << iota
	JPG
	GIF
)

func (mt MediaType) Extension() string {
	switch mt {
	default:
		return ""
	case MP4:
		return "mp4"
	case JPG:
		return "jpg"
	case GIF:
		return "gif"
	}
}

func (mt MediaType) String() string {
	return mt.Extension()
}

var (
	debug = os.Getenv("DEBUG_GIFS_PKG") != ""
)

const (
	importEndpointURL = "https://api.gifs.com/media/import"
)

var (
	defaultConcurrentImportsCount = uint64(10)
)

func debugLogPrintf(fmt_ string, args ...interface{}) {
	if debug {
		log.Printf(fmt_, args...)
	}
}

type Credentials struct {
	APIKey string `json:"api_key,omitempty"`
}

type GIFS struct {
	client *http.Client
	apiKey string

	mu sync.Mutex
}

func (g *GIFS) SetAPIKey(apiKey string) {
	g.mu.Lock()
	g.apiKey = apiKey
	g.mu.Unlock()
}

func New() (*GIFS, error) {
	return new(GIFS), nil
}

func NewWithCredentials(cred *Credentials) (*GIFS, error) {
	g, err := New()
	if err != nil {
		return nil, err
	}
	if cred != nil {
		g.SetAPIKey(cred.APIKey)
	}
	return g, nil
}

type Trim struct {
	Start float64 `json:"start,omitempty"`
	End   float64 `json:"end,omitempty"`
}

type Params struct {
	Title  string   `json:"title,omitempty"`
	URL    string   `json:"source,omitempty"`
	APIKey string   `json:"api_key,omitempty"`
	Tags   []string `json:"tags,omitempty"`
	NSFW   bool     `json:"nsfw,omitempty"`
	Trim   *Trim    `json:"trim,omitempty"`

	// Only set media if you are performing an upload
	media io.Reader

	Attribution map[string]interface{} `json:"attribution,omitempty"`

	callbackURI string `json:"-"`
}

func (p *Params) SetMedia(r io.Reader) error {
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

func (re responseError) Error() string {
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

func (g *GIFS) Upload() (*Response, error) {
	return nil, errUnimplemented
}

// Import is a method with which you'll specify atleast
// an http based URL pointing to media that you'd like
// to import to gifs.com.
func (g *GIFS) Import(param *Params) (*Response, error) {
	if param == nil {
		return nil, ErrNilParamDereference
	}

	bip := &BulkImportParams{
		Params: []*Params{param},
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
// multiple media URLs without having to construct each `Params` object.
func (g *GIFS) ImportSources(sources ...string) ([]*Response, error) {
	if len(sources) < 1 {
		return nil, ErrExpectingAtLeastOneSource
	}

	preparedParams := []*Params{}
	for _, source := range sources {
		preparedParams = append(preparedParams, &Params{URL: source})
	}

	bip := &BulkImportParams{Params: preparedParams}
	return g.ImportBulk(bip)
}

type BulkImportParams struct {
	ConcurrentImports uint
	Params            []*Params
}

func (p *Params) transformToImportBody() ([]byte, error) {
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

func (g *GIFS) doPOSTRequest(uri string, param *Params, headers http.Header) (*http.Response, error) {
	byteSlice, err := param.transformToImportBody()
	if err != nil {
		return nil, err
	}
	debugLogPrintf("byteSlice body %s for param: %+v\n", byteSlice, param)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(byteSlice))
	if err != nil {
		return nil, err
	}

	copyHeaders(headers, req.Header)
	req.Header.Set("Content-Type", "application/json")
	if g.apiKey != "" {
		req.Header.Set("Gifs-Api-Key", g.apiKey)
	}

	return g.createdClient().Do(req)
}

func (g *GIFS) doMultipartUpload() (*http.Response, error) {
	return nil, errUnimplemented
}

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

func (g *GIFS) createdClient() HTTPDoer {
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
	param   *Params
	headers http.Header

	uuid uint64
	g    *GIFS
	typ  jobType
}

func (hj httpRequestJob) Id() interface{} {
	return hj.uuid
}

func (hj httpRequestJob) gifsLiason() *GIFS {
	if hj.g != nil {
		return hj.g
	}

	g, err := New()
	if err != nil {
		panic(err)
	}
	return g
}

func (hj httpRequestJob) Do() (interface{}, error) {
	res, err := hj.gifsLiason().doPOSTRequest(hj.uri, hj.param, hj.headers)
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
		wrapRes := res.(*wrapperResponse)
		if wrapRes != nil {
			finalRes = wrapRes.Success
			if err == nil {
				err = wrapRes.Errors
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
func (g *GIFS) ImportBulk(bip *BulkImportParams) ([]*Response, error) {
	concurrentImports := defaultConcurrentImportsCount
	if bip.ConcurrentImports > 0 {
		concurrentImports = uint64(bip.ConcurrentImports)
	}

	maxResponseId := uint64(len(bip.Params))
	jobsBench := make(chan semalim.Job)
	go func() {
		defer close(jobsBench)
		for i := uint64(0); i < maxResponseId; i++ {
			param := bip.Params[i]
			jobsBench <- httpRequestJob{uri: importEndpointURL, param: param, uuid: uint64(i)}
		}
	}()

	resultsChan := semalim.Run(jobsBench, concurrentImports)
	return categorizeParallelJobResponses(resultsChan, maxResponseId)
}
