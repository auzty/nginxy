package main

import (
    "fmt"
    "log"
    "flag"
    "time"
    "os"
    "github.com/fsouza/go-dockerclient"
)

type Nginx struct{
  ServiceName string
  DomainName string
  ServicePort string
  Ssl string
  SslKey string
  SslCert string
}

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
    panic(fmt.Sprintf("%s: %s", msg, err))
  }
}

var serviceCount = 0
var serviceList []string
func main() {

    filter := make(map[string][]string)
    filter["label"] = append(filter["label"],"nginxy.domain")

    params := flag.String("endpoint","/tmp/docker.sock","Docker Endpoint (socket only)")
    flag.Parse()

    endpoint := "unix://"+*params
    client, err := docker.NewClient(endpoint)
    if err != nil {
      panic(err)
    }

    // get the current list services and service count
    serviceCount,serviceList := getServiceCount(client,filter)

    imgs, err := client.ListServices(docker.ListServicesOptions{Filters: filter})
    if err != nil {
      panic(err)
    }
    for _, img := range imgs {
      if ((serviceCount != 0)&& (img.Spec.Labels["nginxy.domain"] != "")){
        // -- Generate the start configuration 
         name := img.Spec.Name

        // check configuration is exist or not 
         _,filerr := os.Stat("/etc/nginx/conf.d/"+name+".conf")

         if (filerr != nil){ // if didn't have configuration, then create the configuration

           domain := img.Spec.Labels["nginxy.domain"]
           port := img.Spec.Labels["nginxy.port"]
           ssl := img.Spec.Labels["nginxy.ssl"]
           key := img.Spec.Labels["nginxy.ssl.key"]
           cert := img.Spec.Labels["nginxy.ssl.cert"]

           nginxconf := Nginx{ServiceName:name,DomainName:domain,ServicePort:port,Ssl:ssl,SslKey:key,SslCert:cert}
           fmt.Println(nginxconf)
           nginxconf.WriteConf() // write the nginx configuration based on the name
           reloadNginx()
         } else {
          log.Printf("File configuration of %s already exist, skip creating configuration for %s\n",img.Spec.Name)
         }

      }
    }

    log.Println("Jumlah Service : ",serviceCount)

// -- goruoutine start

forever := make(chan bool)

  go func() {
    for 1 < 5 {
      // execution go here

      currentCount,currentList := getServiceCount(client,filter)
//      fmt.Printf("Current list : %v , awal Service List : %v\n",currentList,serviceList)
//      fmt.Printf("Current Count : %v , awal Service Count : %v\n",currentCount,serviceCount)

      // if currentCount > serviceCount do add new services
      if(currentCount >  serviceCount) {
       // new Services, create new configuration

       // get the differences services (check for new services)
       diff := difference(currentList,serviceList)

       for _, nama := range diff {
         svcID := getServiceID(client,filter,nama)

         detail, err := client.InspectService(svcID)

         if err != nil {
          log.Println(err)
         }

         name := detail.Spec.Name
         domain := detail.Spec.Labels["nginxy.domain"]
         port := detail.Spec.Labels["nginxy.port"]
         ssl := detail.Spec.Labels["nginxy.ssl"]
         key := detail.Spec.Labels["nginxy.ssl.key"]
         cert := detail.Spec.Labels["nginxy.ssl.cert"]

         nginxconf := Nginx{ServiceName:name,DomainName:domain,ServicePort:port,Ssl:ssl,SslKey:key,SslCert:cert}
         fmt.Println(nginxconf)
         nginxconf.WriteConf() // write the nginx configuration based on the name
         reloadNginx()

       }

        serviceList = currentList
        serviceCount++
        // increment the serviceCount
      }

      if(currentCount < serviceCount) {
      //  fmt.Printf("masuk hapus current list : %v , Service List : %v",currentList,serviceList)
        // delete the configuration because the service is deleted

        // get the differences
        diff := difference(currentList,serviceList)

        for _, nama := range diff {
          path := "/etc/nginx/conf.d/"+nama+".conf"
          err := os.Remove(path)
          if err != nil {
            log.Println(err)
          }else{
            serviceCount--
            // how to delete the element ??
            serviceList = currentList
            reloadNginx()
            log.Println("Removing service ",nama)
          }
        }

      }

      time.Sleep(1 * time.Second)
      // end
    }
  }()

log.Printf(" [*] Waiting for Docker Event(s). CTRL+C to exit\n")
<-forever

// -- goroutine end


}


// to get current service list
func getServiceCount(client *docker.Client,filter map[string][]string) (int ,[]string) {

  var listSvc []string
  localCount := 0
  imgs, err := client.ListServices(docker.ListServicesOptions{Filters: filter})
    if err != nil {
        panic(err)
    }
    for _,img := range imgs {
      listSvc = append(listSvc,img.Spec.Name)
      localCount++
    }
 return localCount,listSvc
}


// to get list of differences between 2 slices
func difference(slice1 []string, slice2 []string) []string {
  var diff []string

  // Loop two times, first to find slice1 strings not in slice2,
  // second loop to find slice2 strings not in slice1
  for i := 0; i < 2; i++ {
    for _, s1 := range slice1 {
      found := false
      for _, s2 := range slice2 {
        if s1 == s2 {
          found = true
          break
        }
      }
      // String not found. We add it to return slice
      if !found {
        diff = append(diff, s1)
      }
    }
    // Swap the slices, only if it was the first loop
    if i == 0 {
      slice1, slice2 = slice2, slice1
    }
  }

  return diff
}

// find the current Service ID's 
func getServiceID( client *docker.Client,filter map[string][]string,name string) string {
  imgs, err := client.ListServices(docker.ListServicesOptions{Filters: filter})
    if err != nil {
        panic(err)
    }
    for _,img := range imgs {
      if(img.Spec.Name == name){
        return img.ID
        break
      }
    }
    return ""
}
