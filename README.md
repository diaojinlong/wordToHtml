#Word转Html文本服务

```shell
docker run -itd --name wordToHtml -p 8083:8083 diaojinlong/word_to_html:v1.0.0
```

# 使用方法
post

http://127.0.0.1:8083/convert
file=1.docx

```shell
POST /convert HTTP/1.1
Host: 127.0.0.1:8083
Content-Length: 262
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW

------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="file"; filename="/E:/docker/LibreOffice/1.docx"
Content-Type: application/vnd.openxmlformats-officedocument.wordprocessingml.document

(data)
------WebKitFormBoundary7MA4YWxkTrZu0gW--

```