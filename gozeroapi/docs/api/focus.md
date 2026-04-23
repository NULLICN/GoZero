### 1. "获取焦点图列表"

1. route definition

- Url: /api/focus
- Method: GET
- Request: `-`
- Response: `CommonResponse`

2. request definition



3. response definition



```golang
type CommonResponse struct {
	Success bool `json:"success"`
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}
```

### 2. "通过id获取一个焦点图"

1. route definition

- Url: /api/focus/:id
- Method: GET
- Request: `FocusRequestByPath`
- Response: `CommonResponse`

2. request definition



```golang
type FocusRequestByPath struct {
	Id string `path:"id"` // 动态路由传值 /:id 术语：路径参数（Path Parameter）或 URL 参数（URL Parameter）
}
```


3. response definition



```golang
type CommonResponse struct {
	Success bool `json:"success"`
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}
```

### 3. 通过请求体获取焦点图

1. route definition

- Url: /api/focus/body
- Method: POST
- Request: `FocusRequestByBody`
- Response: `CommonResponse`

2. request definition



```golang
type FocusRequestByBody struct {
	Id string `form:"id"` // 表单传值
}
```


3. response definition



```golang
type CommonResponse struct {
	Success bool `json:"success"`
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}
```

### 4. "通过查询参数获取焦点图"

1. route definition

- Url: /api/focus/query
- Method: GET
- Request: `FocusRequestByQuery`
- Response: `CommonResponse`

2. request definition



```golang
type FocusRequestByQuery struct {
	Id string `form:"id,default=&#39;1145&#39;"` // 查询传值 ?id=123 术语：查询参数（Query Parameter）或 查询字符串（Query String）
}
```


3. response definition



```golang
type CommonResponse struct {
	Success bool `json:"success"`
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}
```

