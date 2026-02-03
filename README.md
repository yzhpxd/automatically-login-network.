# automatically-login-network.
路由自动登录外网

调试了小米ax5jdc，由于公司外网每天都要登录，go了个小程序，让路由器自动登录。关键是提取登录密码，登录界面按f12，找到<img width="1142" height="484" alt="image" src="https://github.com/user-attachments/assets/923c10f3-ca4a-40d8-bcac-8e3cf78c8aa0" />，复制代码curl 'http://192.168.3.1/ac_portal/login.php' \
  -H 'Referer: http://192.168.3.1/ac_portal/default/pc.html?template=default&tabs=pwd&vlanid=0&_ID_=0&switch_url=&url=http://192.168.3.1/homepage/index.html&controller_type=&mac=8c-de-f9-1c-0a-dd' \
  -H 'X-Requested-With: XMLHttpRequest' \
  -H 'User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36' \
  -H 'Accept: */*' \
  -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' \
  --data-raw 'opr=pwdLogin&userName=aaaaaa&pwd=bbbbbbb&auth_tag=xxxxxx&rememberPwd=0' \
  --insecure
得到登录密码和对应的时间，这两个是对应的，相当于配对
