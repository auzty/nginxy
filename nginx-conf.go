package main

import(
  "fmt"
  "log"
  "io/ioutil"
  "text/template"
  "os"
  "strings"
  "strconv"
  "time"
  "syscall"
)

func (org *Nginx) WriteConf() error {
  template_path := "/etc/nginx/templating/conf.tmpl"
  buf, err := ioutil.ReadFile(template_path)
  failOnError(err,"template file not found")

  text := string(buf)
  conf_file := fmt.Sprintf("%s.conf", org.ServiceName)
  fd, err := os.Create(fmt.Sprintf("/etc/nginx/conf.d/%s",conf_file))
  failOnError(err,"File Create Error")

  defer func(){
    fd.Close()
  }()

  tmpl, err := template.New("conf").Parse(text)
  failOnError(err,"template parsing error")

  tmpl.Execute(fd,org)
  return nil
}

func reloadNginx(){
  pidloc := "/var/run/nginx.pid"

  nginxpid, nginxerr := ioutil.ReadFile(pidloc)
  failOnError(nginxerr,"PID not found / nginx is not running")

  // convert []bytes become INT
  pid_str := strings.Replace(string(nginxpid),"\n","",-1)
  pid, err := strconv.Atoi(pid_str)

  log.Println("Reloading Nginx Configuration.....")
  time.Sleep(2 * time.Second)

  if nginxerr == nil {
    err =  syscall.Kill(pid,syscall.SIGHUP)
    failOnError(err,"Failed send SIGHUP to nginx")
  }

}
