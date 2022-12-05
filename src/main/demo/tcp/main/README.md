1.  上传文件
scp ./server.go pengchao@xxx.xxx.x.x:~/Uploads
2.  移动到指定目录
mv server.go /usr/local/project
3. go build (需要把localhost改成0.0.0.0, 这样才能外网请求)
4. ./server 启动


注意：
需要打开linux的防火墙
查看：systemctl status firewalld
关闭：systemctl stop firewalld