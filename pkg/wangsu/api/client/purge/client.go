// This file is auto-generated, don't edit it. Thanks.
package purge

import (
	"github.com/alibabacloud-go/tea/tea"
)

type Paths struct {
}

func (s Paths) String() string {
	return tea.Prettify(s)
}

func (s Paths) GoString() string {
	return s.String()
}

type Parameters struct {
}

func (s Parameters) String() string {
	return tea.Prettify(s)
}

func (s Parameters) GoString() string {
	return s.String()
}

type RequestHeader struct {
}

func (s RequestHeader) String() string {
	return tea.Prettify(s)
}

func (s RequestHeader) GoString() string {
	return s.String()
}

type CreateAPurgeRequestRequest struct {
	// {"en" : "A description of the purge request.", "zh_CN": "刷新请求的说明。"}
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// {"en" : "URLs of files to purge.  File URLs should not contain the asterisk character, '*'.   If a directory or filename in a URL includes a percent character, '%', be sure to encode it. A URL can be up to 2048 characters.", "zh_CN": "要刷新的文件的URL。URL不能包含星号字符'*'。如果URL中的目录或文件名包含'%'等特殊符号，需要先进行URL编码。每个URL长度不能超过2048个字符。"}
	FileUrls []*string `json:"fileUrls,omitempty" xml:"fileUrls,omitempty" type:"Repeated"`
	// {"en" : "If a file's cache key depends on request headers, you can specify the header values that are applicable to purge one version of the cached file. The same set of header values will apply to all entries in fileUrls.", "zh_CN": "如果文件的缓存键与请求头相关，则可以指定请求头和值来刷新相应的缓存文件。此处指定的请求头和值将应用于fileUrls中的所有条目。"}
	FileHeaders []*CreateAPurgeRequestRequestFileHeaders `json:"fileHeaders,omitempty" xml:"fileHeaders,omitempty" type:"Repeated"`
	// {"en" : "<= 20 items
	// URLs to purge. URLs must begin with http:// or https:// and can be up to 2048 characters. Use the '*' character to purge multiple files or directories. If a URL has multiple sets of asterisk characters, only the last '*' or '**' will be treated as a wildcard. Other instances of '*' earlier in the URL will be treated as the literal character '*'.
	// <table><tr><th>Example</th><th>Description</th></tr><tr><td>http://test.domain2.com/mydir</td><td>Purge all variations of a single directory, but not its subdirectories or files. Variations may exist if custom cache keys are used.</td></tr><tr><td>http://test.domain2.com/mydir/**</td><td>Purge all files and subdirectories whose cache key begins with http://test.domain2.com/mydir/.</td></tr><tr><td>http://test.domain2.com/mydir/*</td><td>Purge all files, but not subdirectories, within a directory.</td></tr><tr><td>http://test.domain2.com/mydir/*.jpg</td><td>Purge all cache entries ending with the .jpg file extension. Subdirectories of http://test.domain2.com/mydir/ are not purged. </td></tr><tr><td>http://test.domain2.com/mydir/a*</td><td>Purge all files, but not subdirectories, that start with the letter 'a'.</td></tr><tr><td>http://test.domain2.com/mydir/a**</td><td>Purge all files and subdirectories that start with the letter 'a'.</td></tr><tr><td>http://test.domain2.com/mydir/a.jpg</td><td>Purge all variations of 'a.jpg'. Variations may exist if custom cache keys are used.</td></tr><tr><td>http://test.domain2.com/my**jpg</td><td>Purge all entries whose cache key begins with http://test.domain2.com/my and ends with the suffix jpg. The '**' can match anything in the path including additional subdirectories. For example, http://test.domain2.com/mydirectory/picture.jpg would be purged.</td></tr></table>
	// If a directory or filename in a URL includes a percent character, '%', be sure to encode it.", "zh_CN": "<= 20 条目
	// 要刷新的目录的URL。URL必须以http:// 或者 https://开头，每条URL最多只能包含2048个字符。 在URL中使用'*'字符可以匹配多个文件或目录。如果一条URL中带有多组'*'，则只有最后一个'*'或'**'会被当成通配符来进行匹配，其它的'*'只会被当成普通字符。
	// <table><tr><th>示例</th><th>描述</th></tr><tr><td>http://test.domain2.com/mydir</td><td>刷新单个目录的所有变体，但不包括其子目录或文件。当您自定义了缓存键时，则可能存在变体。</td></tr><tr><td>http://test.domain2.com/mydir/**</td><td>刷新缓存键以http://test.domain2.com/mydir/开头的所有文件和子目录。</td></tr><tr><td>http://test.domain2.com/mydir/*</td><td>刷新目录中的所有文件，但不包括子目录。</td></tr><tr><td>http://test.domain2.com/mydir/*.jpg</td><td>刷新所有以.jpg文件扩展名结尾的缓存，但不会刷新http://test.domain2.com/mydir/ 的子目录。 </td></tr><tr><td>http://test.domain2.com/mydir/a*</td><td>刷新以字母'a'开头的所有文件，但不包括子目录。</td></tr><tr><td>http://test.domain2.com/mydir/a**</td><td>刷新以字母'a'开头的所有文件和子目录。</td></tr><tr><td>http://test.domain2.com/mydir/a.jpg</td><td>刷新'a.jpg'文件的所有变体。当您自定义了缓存键时，则可能存在变体。</td></tr><tr><td>http://test.domain2.com/my**jpg</td><td>刷新缓存键以 http://test.domain2.com/my 开头并以后缀 jpg 结尾的所有条目。'**'可以匹配路径中的任何内容，包括其他子目录。例如，http://test.domain2.com/mydirectory/picture.jpg 将被刷新。</td></tr></table>
	// 如果URL中的目录或文件名包含百分号'%'等特殊符号时，请确保先进行URL编码。"}
	DirUrls []*string `json:"dirUrls,omitempty" xml:"dirUrls,omitempty" type:"Repeated"`
	// {"en" : "<= 2 items
	// Regular expression patterns used to match the cache key. Each must begin with the following format:
	//  {scheme}://{hostname}/. {scheme} can be http, https, or any, which matches any scheme.
	// Example:
	// https://test.domain.com/my.*\.(jpg|png)\?q=
	// <br/>
	// For performance considerations, the following restrictions apply:
	// The regular expression pattern following the hostname can be up to 126 characters.
	//
	// It can consist of up to two unlimited quantifiers ('*', '+', or ',}').
	// The upper limit on a quantifier cannot be more than 59, for example, {1,59}", "zh_CN": "<= 2 条目
	// 用于匹配缓存键的正则表达式。
	// 每个表达式必须以
	// {协议}://{域名}/ 格式开头。其中，{协议} 可以是 http, https，或any（表示不限协议）。
	// 示例：
	// https://test.domain.com/my.*\.(jpg|png)\?q=
	// <br/>
	// 出于性能考虑，使用正则表达式有以下限制：
	//
	// 在域名后面的正则表达式最多只能包含126个字符。
	// 最多只能包含两个限定符('*'、'+'或',}')。
	// 限定符的上限不能超过59，例如{1,59}"}
	RegexPatterns []*string `json:"regexPatterns,omitempty" xml:"regexPatterns,omitempty" type:"Repeated"`
	// {"en" : "Enum: delete invalidate
	// Default: invalidate
	// This controls whether cached files and directories should be removed altogether from the CDN Pro servers (delete) or flagged as invalid (invalidate).", "zh_CN": "取值范围: delete, invalidate
	// 默认值: invalidate
	// 指定刷新类型，包括完全删除(delete)和标记为无效(invalidate)。"}
	Action *string `json:"action,omitempty" xml:"action,omitempty"`
	// {"en" : "Enum: staging production
	// Specify if the purge request applies to the staging or production environment.", "zh_CN": "取值范围: staging, production
	// 指定刷新请求应用于演练环境还是生产环境。"}
	Target *string `json:"target,omitempty" xml:"target,omitempty" require:"true"`
	// {"en" : "ID of a webhook to call when the purge task completes.", "zh_CN": "刷新任务完成时要调用的webhook的ID。"}
	Webhook *string `json:"webhook,omitempty" xml:"webhook,omitempty"`
}

func (s CreateAPurgeRequestRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateAPurgeRequestRequest) GoString() string {
	return s.String()
}

func (s *CreateAPurgeRequestRequest) SetName(v string) *CreateAPurgeRequestRequest {
	s.Name = &v
	return s
}

func (s *CreateAPurgeRequestRequest) SetFileUrls(v []*string) *CreateAPurgeRequestRequest {
	s.FileUrls = v
	return s
}

func (s *CreateAPurgeRequestRequest) SetFileHeaders(v []*CreateAPurgeRequestRequestFileHeaders) *CreateAPurgeRequestRequest {
	s.FileHeaders = v
	return s
}

func (s *CreateAPurgeRequestRequest) SetDirUrls(v []*string) *CreateAPurgeRequestRequest {
	s.DirUrls = v
	return s
}

func (s *CreateAPurgeRequestRequest) SetRegexPatterns(v []*string) *CreateAPurgeRequestRequest {
	s.RegexPatterns = v
	return s
}

func (s *CreateAPurgeRequestRequest) SetAction(v string) *CreateAPurgeRequestRequest {
	s.Action = &v
	return s
}

func (s *CreateAPurgeRequestRequest) SetTarget(v string) *CreateAPurgeRequestRequest {
	s.Target = &v
	return s
}

func (s *CreateAPurgeRequestRequest) SetWebhook(v string) *CreateAPurgeRequestRequest {
	s.Webhook = &v
	return s
}

type CreateAPurgeRequestRequestFileHeaders struct {
	// {"en" : "HTTP header name.", "zh_CN": "HTTP 头部名称"}
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// {"en" : "Value of an HTTP header.", "zh_CN": "HTTP 头部的值"}
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s CreateAPurgeRequestRequestFileHeaders) String() string {
	return tea.Prettify(s)
}

func (s CreateAPurgeRequestRequestFileHeaders) GoString() string {
	return s.String()
}

func (s *CreateAPurgeRequestRequestFileHeaders) SetName(v string) *CreateAPurgeRequestRequestFileHeaders {
	s.Name = &v
	return s
}

func (s *CreateAPurgeRequestRequestFileHeaders) SetValue(v string) *CreateAPurgeRequestRequestFileHeaders {
	s.Value = &v
	return s
}

type CreateAPurgeRequestResponse struct {
}

func (s CreateAPurgeRequestResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateAPurgeRequestResponse) GoString() string {
	return s.String()
}

type ResponseHeader struct {
	// {"en":"The Location header is a URL representing the new purge request, for example, <code>Location: http://ngapi.cdnetworks.com/cdn/purges/e91e8674-c2c5-4440-a1de-8b2ea99293dd</code>.", "zh_CN":"通过Location响应头返回新建的刷新任务的URL。URL中包含刷新任务的ID，可使用该ID调用'查询刷新任务详情'接口来查看刷新任务详情。URL示例：<code>Location: http://open.chinanetcenter.com/cdn/purges/5dca2205f9e9cc0001df7b33"}
	Location *string `json:"Location,omitempty" xml:"Location,omitempty" require:"true"`
}

func (s ResponseHeader) String() string {
	return tea.Prettify(s)
}

func (s ResponseHeader) GoString() string {
	return s.String()
}

func (s *ResponseHeader) SetLocation(v string) *ResponseHeader {
	s.Location = &v
	return s
}
