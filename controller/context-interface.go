package controller

import (
	"io"
	"mime/multipart"
	"time"

	"context"
)

// ErrorType is an unsigned 64-bit error code as defined in the gin spec.
type ErrorType = uint64

const (
	// ErrorTypeBind is used when Context.Bind() fails.
	ErrorTypeBind ErrorType = 1 << 63
	// ErrorTypeRender is used when Context.Render() fails.
	ErrorTypeRender ErrorType = 1 << 62
	// ErrorTypePrivate indicates a private error.
	ErrorTypePrivate ErrorType = 1 << 0
	// ErrorTypePublic indicates a public error.
	ErrorTypePublic ErrorType = 1 << 1
	// ErrorTypeAny indicates any other error.
	ErrorTypeAny ErrorType = 1<<64 - 1
	// ErrorTypeNu indicates any other error.
	ErrorTypeNu = 2
)

// Error represents a error's specification.
type Error interface {
	GetError() error
	GetType() ErrorType
	GetMeta() interface{}
}

// MContext is the interface of Control Layer
type MContext interface {
	// MContext must be a context.Context first
	context.Context
	// Copy returns a copy of the current context that can be safely used outside the request's scope.
	// This has to be used when the context has to be passed to a goroutine.
	Copy() MContext
	// HandlerName returns the main handler's name. For example if the handler is "handleGetUsers()",
	// this function will return "main.handleGetUsers".
	HandlerName() string
	// HandlerNames returns a list of all registered handlers for this context in descending order,
	// following the semantics of HandlerName()
	HandlerNames() []string
	// FullPath returns a matched route full path. For not found routes
	// returns an empty string.
	//     router.GET("/user/:id", func(c mcontext.MContext) {
	//         c.FullPath() == "/user/:id" // true
	//     })
	FullPath() string
	// Next should be used only inside middleware.
	// It executes the pending handlers in the chain inside the calling handler.
	// See example in GitHub.
	Next()
	// IsAborted returns true if the current context was aborted.
	IsAborted() bool
	// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
	// Let's say you have an authorization middleware that validates that the current request is authorized.
	// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
	// for this request are not called.
	Abort()
	// AbortWithStatus calls `Abort()` and writes the headers with the specified status code.
	// For example, a failed attempt to authenticate a request could use: context.AbortWithStatus(401).
	AbortWithStatus(code int)
	// AbortWithStatusJSON calls `Abort()` and then `JSON` internally.
	// This method stops the chain, writes the status code and return a JSON body.
	// It also sets the Content-Type as "application/json".
	AbortWithStatusJSON(code int, jsonObj interface{})
	// AbortWithError calls `AbortWithStatus()` and `Error()` internally.
	// This method stops the chain, writes the status code and pushes the specified error to `c.Errors`.
	// See Context.Error() for more details.
	AbortWithError(code int, err error) Error
	// Error attaches an error to the current context. The error is pushed to a list of errors.
	// It's a good idea to call Error for each error that occurred during the resolution of a request.
	// A middleware can be used to collect all the errors and push them to a database together,
	// print a log, or append it in the HTTP response.
	// Error will panic if err is nil.
	Error(err error) Error
	// Set is used to store a new key/value pair exclusively for this context.
	// It also lazy initializes  c.Keys if it was not used previously.
	Set(key string, value interface{})
	// Get returns the value for the given key, ie: (value, true).
	// If the value does not exists it returns (nil, false)
	Get(key string) (value interface{}, exists bool)
	// MustGet returns the value for the given key if it exists, otherwise it panics.
	MustGet(key string) interface{}
	// GetString returns the value associated with the key as a string.
	GetString(key string) (s string)
	// GetBool returns the value associated with the key as a boolean.
	GetBool(key string) (b bool)
	// GetInt returns the value associated with the key as an integer.
	GetInt(key string) (i int)
	// GetInt64 returns the value associated with the key as an integer.
	GetInt64(key string) (i64 int64)
	// GetFloat64 returns the value associated with the key as a float64.
	GetFloat64(key string) (f64 float64)
	// GetTime returns the value associated with the key as time.
	GetTime(key string) (t time.Time)
	// GetDuration returns the value associated with the key as a duration.
	GetDuration(key string) (d time.Duration)
	// GetStringSlice returns the value associated with the key as a slice of strings.
	GetStringSlice(key string) (ss []string)
	// GetStringMap returns the value associated with the key as a map of interfaces.
	GetStringMap(key string) (sm map[string]interface{})
	// GetStringMapString returns the value associated with the key as a map of strings.
	GetStringMapString(key string) (sms map[string]string)
	// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
	GetStringMapStringSlice(key string) (smss map[string][]string)
	// Param returns the value of the URL param.
	// It is a shortcut for c.Params.ByName(key)
	//     router.GET("/user/:id", func(c mcontext.MContext) {
	//         // a GET request to /user/john
	//         id := c.Param("id") // id == "john"
	//     })
	Param(key string) string
	// Query returns the keyed url query value if it exists,
	// otherwise it returns an empty string `("")`.
	// It is shortcut for `c.Request.URL.Query().Get(key)`
	//     GET /path?id=1234&name=Manu&value=
	// 	   c.Query("id") == "1234"
	// 	   c.Query("name") == "Manu"
	// 	   c.Query("value") == ""
	// 	   c.Query("wtf") == ""
	Query(key string) string
	// DefaultQuery returns the keyed url query value if it exists,
	// otherwise it returns the specified defaultValue string.
	// See: Query() and GetQuery() for further information.
	//     GET /?name=Manu&lastname=
	//     c.DefaultQuery("name", "unknown") == "Manu"
	//     c.DefaultQuery("id", "none") == "none"
	//     c.DefaultQuery("lastname", "none") == ""
	DefaultQuery(key, defaultValue string) string
	// GetQuery is like Query(), it returns the keyed url query value
	// if it exists `(value, true)` (even when the value is an empty string),
	// otherwise it returns `("", false)`.
	// It is shortcut for `c.Request.URL.Query().Get(key)`
	//     GET /?name=Manu&lastname=
	//     ("Manu", true) == c.GetQuery("name")
	//     ("", false) == c.GetQuery("id")
	//     ("", true) == c.GetQuery("lastname")
	GetQuery(key string) (string, bool)
	// QueryArray returns a slice of strings for a given query key.
	// The length of the slice depends on the number of params with the given key.
	QueryArray(key string) []string
	// GetQueryArray returns a slice of strings for a given query key, plus
	// a boolean value whether at least one value exists for the given key.
	GetQueryArray(key string) ([]string, bool)
	// QueryMap returns a map for a given query key.
	QueryMap(key string) map[string]string
	// GetQueryMap returns a map for a given query key, plus a boolean value
	// whether at least one value exists for the given key.
	GetQueryMap(key string) (map[string]string, bool)
	// PostForm returns the specified key from a POST urlencoded form or multipart form
	// when it exists, otherwise it returns an empty string `("")`.
	PostForm(key string) string
	// DefaultPostForm returns the specified key from a POST urlencoded form or multipart form
	// when it exists, otherwise it returns the specified defaultValue string.
	// See: PostForm() and GetPostForm() for further information.
	DefaultPostForm(key, defaultValue string) string
	// GetPostForm is like PostForm(key). It returns the specified key from a POST urlencoded
	// form or multipart form when it exists `(value, true)` (even when the value is an empty string),
	// otherwise it returns ("", false).
	// For example, during a PATCH request to update the user's email:
	//     email=mail@example.com  -->  ("mail@example.com", true) := GetPostForm("email") // set email to "mail@example.com"
	// 	   email=                  -->  ("", true) := GetPostForm("email") // set email to ""
	//                             -->  ("", false) := GetPostForm("email") // do nothing with email
	GetPostForm(key string) (string, bool)
	// PostFormArray returns a slice of strings for a given form key.
	// The length of the slice depends on the number of params with the given key.
	PostFormArray(key string) []string
	// GetPostFormArray returns a slice of strings for a given form key, plus
	// a boolean value whether at least one value exists for the given key.
	GetPostFormArray(key string) ([]string, bool)
	// PostFormMap returns a map for a given form key.
	PostFormMap(key string) map[string]string
	// GetPostFormMap returns a map for a given form key, plus a boolean value
	// whether at least one value exists for the given key.
	GetPostFormMap(key string) (map[string]string, bool)
	// FormFile returns the first file for the provided form key.
	FormFile(name string) (*multipart.FileHeader, error)
	// MultipartForm is the parsed multipart form, including file uploads.
	MultipartForm() (*multipart.Form, error)
	// SaveUploadedFile uploads the form file to specific dst.
	SaveUploadedFile(file *multipart.FileHeader, dst string) error
	// Bind checks the Content-Type to select a binding engine automatically,
	// Depending the "Content-Type" header different bindings are used:
	//     "application/json" --> JSON binding
	//     "application/xml"  --> XML binding
	// otherwise --> returns an error.
	// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
	// It decodes the json payload into the struct specified as a pointer.
	// It writes a 400 error and sets Content-Type header "text/plain" in the response if input is not valid.
	Bind(obj interface{}) error
	// BindJSON is a shortcut for c.MustBindWith(obj, binding.JSON).
	BindJSON(obj interface{}) error
	// BindXML is a shortcut for c.MustBindWith(obj, binding.BindXML).
	BindXML(obj interface{}) error
	// BindQuery is a shortcut for c.MustBindWith(obj, binding.Query).
	BindQuery(obj interface{}) error
	// BindYAML is a shortcut for c.MustBindWith(obj, binding.YAML).
	BindYAML(obj interface{}) error
	// BindHeader is a shortcut for c.MustBindWith(obj, binding.Header).
	BindHeader(obj interface{}) error
	// BindUri binds the passed struct pointer using binding.Uri.
	// It will abort the request with HTTP 400 if any error occurs.
	BindUri(obj interface{}) error
	// ShouldBind checks the Content-Type to select a binding engine automatically,
	// Depending the "Content-Type" header different bindings are used:
	//     "application/json" --> JSON binding
	//     "application/xml"  --> XML binding
	// otherwise --> returns an error
	// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
	// It decodes the json payload into the struct specified as a pointer.
	// Like c.Bind() but this method does not set the response status code to 400 and abort if the json is not valid.
	ShouldBind(obj interface{}) error
	// ShouldBindJSON is a shortcut for c.ShouldBindWith(obj, binding.JSON).
	ShouldBindJSON(obj interface{}) error
	// ShouldBindXML is a shortcut for c.ShouldBindWith(obj, binding.XML).
	ShouldBindXML(obj interface{}) error
	// ShouldBindQuery is a shortcut for c.ShouldBindWith(obj, binding.Query).
	ShouldBindQuery(obj interface{}) error
	// ShouldBindYAML is a shortcut for c.ShouldBindWith(obj, binding.YAML).
	ShouldBindYAML(obj interface{}) error
	// ShouldBindHeader is a shortcut for c.ShouldBindWith(obj, binding.Header).
	ShouldBindHeader(obj interface{}) error
	// ShouldBindUri binds the passed struct pointer using the specified binding engine.
	ShouldBindUri(obj interface{}) error
	// ClientIP implements a best effort algorithm to return the real client IP, it parses
	// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
	// Use X-Forwarded-For before X-Real-Ip as nginx uses X-Real-Ip with the proxy's IP.
	ClientIP() string
	// ContentType returns the Content-Type header of the request.
	ContentType() string
	// IsWebsocket returns true if the request headers indicate that a websocket
	// handshake is being initiated by the client.
	IsWebsocket() bool
	// Status sets the HTTP response code.
	Status(code int)
	// Header is a intelligent shortcut for c.Writer.Header().Set(key, value).
	// It writes a header in the response.
	// If value == "", this method removes the header `c.Writer.Header().Del(key)`
	Header(key, value string)
	// GetHeader returns value from request headers.
	GetHeader(key string) string
	// GetRawData return stream data.
	GetRawData() ([]byte, error)
	// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
	// The provided cookie must have a valid Name. Invalid cookies may be
	// silently dropped.
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	// Cookie returns the named cookie provided in the request or
	// ErrNoCookie if not found. And return the named cookie is unescaped.
	// If multiple cookies match the given name, only one cookie will
	// be returned.
	Cookie(name string) (string, error)
	// HTML renders the HTTP template specified by its file name.
	// It also updates the HTTP code and sets the Content-Type as "text/html".
	// See http://golang.org/doc/articles/wiki/
	HTML(code int, name string, obj interface{})
	// IndentedJSON serializes the given struct as pretty JSON (indented + endlines) into the response body.
	// It also sets the Content-Type as "application/json".
	// WARNING: we recommend to use this only for development purposes since printing pretty JSON is
	// more CPU and bandwidth consuming. Use Context.JSON() instead.
	IndentedJSON(code int, obj interface{})
	// SecureJSON serializes the given struct as Secure JSON into the response body.
	// Default prepends "while(1)," to response body if the given struct is array values.
	// It also sets the Content-Type as "application/json".
	SecureJSON(code int, obj interface{})
	// JSONP serializes the given struct as JSON into the response body.
	// It add padding to response body to request data from a server residing in a different domain than the client.
	// It also sets the Content-Type as "application/javascript".
	JSONP(code int, obj interface{})
	// JSON serializes the given struct as JSON into the response body.
	// It also sets the Content-Type as "application/json".
	JSON(code int, obj interface{})
	// AsciiJSON serializes the given struct as JSON into the response body with unicode to ASCII string.
	// It also sets the Content-Type as "application/json".
	AsciiJSON(code int, obj interface{})
	// PureJSON serializes the given struct as JSON into the response body.
	// PureJSON, unlike JSON, does not replace special html characters with their unicode entities.
	PureJSON(code int, obj interface{})
	// XML serializes the given struct as XML into the response body.
	// It also sets the Content-Type as "application/xml".
	XML(code int, obj interface{})
	// YAML serializes the given struct as YAML into the response body.
	YAML(code int, obj interface{})
	// ProtoBuf serializes the given struct as ProtoBuf into the response body.
	ProtoBuf(code int, obj interface{})
	// String writes the given string into the response body.
	String(code int, format string, values ...interface{})
	// Redirect returns a HTTP redirect to the specific location.
	Redirect(code int, location string)
	// Data writes some data into the body stream and updates the HTTP code.
	Data(code int, contentType string, data []byte)
	// DataFromReader writes the specified reader into the body stream and updates the HTTP code.
	DataFromReader(code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string)
	// File writes the specified file into the body stream in a efficient way.
	File(filepath string)
	// FileAttachment writes the specified file into the body stream in an efficient way
	// On the client side, the file will typically be downloaded with the given filename
	FileAttachment(filepath, filename string)
	// SSEvent writes a Server-Sent Event into the body stream.
	SSEvent(name string, message interface{})
	// Stream sends a streaming response and returns a boolean
	// indicates "Is client disconnected in middle of stream"
	Stream(step func(w io.Writer) bool) bool
	//// Negotiate calls different Render according acceptable Accept format.
	//Negotiate(code int, config Negotiate)
	//// NegotiateFormat returns an acceptable Accept format.
	//NegotiateFormat(offered ...string) string
	// SetAccepted sets Accept header data.
	SetAccepted(formats ...string)

	GetMeta() interface{}
	SetMeta(meta interface{})
}
