# 1. 粘包的解包方式
由于tcp协议是按照流的方式传输，接收端无法知晓发送端发送的数据包的方式和大小。解包主要有三种方式。
## 固定包长
发送方每次发送的包大小固定，接收方收到数据按照固定大小进行处理。针对一些数据长度不变的简单场景，例如传感器数据采集可以采用。
## 指定字符
指定一个特殊的字符做为包的结尾，只要接收方检测到此特殊字符就表示数据包接收完成。像smtp，以crlf结尾。如果传输的数据内容包含特殊字符，要对内容中的特殊字符做转义处理。
## 包头包体
包头固定，可以描述数据包的结构。根据结构描述，再确认包体边界。例如mqtt协议，包头描述协议版本和一些控制信息，payload大小等，然后完整拼接数据包。又如像http，header内以指定字符来分隔消息描述，使用Content-Length确认body大小。

# 2. 协议解析
> buf.read()
## 处理逻辑
1. 读取网络头4个字节的pack size，获取整个goim包大小；
2. 根据读取的pack size获取对应的整个包；
3. 解整个包，获取Operation、SequenceId和Body
