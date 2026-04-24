<h1>goctl命令</h1>
<h3>1.goctl api new 创建一个api类型服务</h3>
<p>goctl api new firstdemo</p>
<h3>2.goctl api go 根据api文件生成go代码</h3>
<p>goctl api go --api firstdemo.api --dir .</p>
<h3>3.goctl api --o 生成api文件</h3>
<p>goctl api --o gozero.api</p>
<h3>4.goctl model mysql ddl 生成SQL模型代码</h3>
<p>goctl model mysql ddl --src ./users.sql --dir .</p>