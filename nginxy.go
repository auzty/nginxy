package main

import (
    "fmt"
    "log"
    "flag"
//    "time"
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
//    imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
//    test := make(map[string][]string)
//    test["label"] = append(test["label"],"nginxy.ssl")
//    imgs, err := client.ListServices(docker.ListServicesOptions{Filters: test })
    imgs, err := client.ListServices(docker.ListServicesOptions{Filters: filter})
    if err != nil {
      panic(err)
    }
    for _, img := range imgs {
        serviceCount++ //increment the total service
//      fmt.Println("Services : ", img)
//      fmt.Println("Services name : ", img.Spec.Name)
//      fmt.Println("Services label : ", img.Spec.Labels)




      //--testing split label --
//      fmt.Println("Labelnya : ",img.Spec.Labels["nginxy.domain"])
//      fmt.Println("Labelnya : ",img.Spec.Labels["nginxy.ssl"])
//      fmt.Println("Labelnya : ",img.Spec.Labels["nginxy.ssl.cert"])
//      fmt.Println("Labelnya : ",img.Spec.Labels["nginxy.ssl.key"])
//        fmt.Println("ID: ", img.ID)
//        fmt.Println("RepoTags: ", img.RepoTags)
//        fmt.Println("Created: ", img.Created)
//        fmt.Println("Size: ", img.Size)
//        fmt.Println("VirtualSize: ", img.VirtualSize)
//        fmt.Println("ParentId: ", img.ParentID)
      if ((serviceCount != 0)&& (img.Spec.Labels["nginxy.domain"] != "")){
        // -- Generate the start configuration 
         name := img.Spec.Name
         domain := img.Spec.Labels["nginxy.domain"]
         port := img.Spec.Labels["nginxy.port"]
         ssl := img.Spec.Labels["nginxy.ssl"]
         key := img.Spec.Labels["nginxy.ssl.key"]
         cert := img.Spec.Labels["nginxy.ssl.cert"]

         nginxconf := Nginx{ServiceName:name,DomainName:domain,ServicePort:port,Ssl:ssl,SslKey:key,SslCert:cert}
         fmt.Println(nginxconf)
         nginxconf.WriteConf() // write the nginx configuration based on the name
         reloadNginx()

      }
    }

    log.Println("Jumlah Service : ",serviceCount)





if err != nil {
    log.Fatal(err)
}

listener := make(chan *docker.APIEvents)
err = client.AddEventListener(listener)
if err != nil {
    log.Fatal(err)
}

defer func() {

    err = client.RemoveEventListener(listener)
    if err != nil {
        log.Fatal(err)
    }

}()

//timeout := time.After(100 * time.Second)

for {
    select {
    case msg := <-listener:


      //if starting container (after docker service create) , will trigger this
      action := msg.Action
      status := msg.Type

/*
      // testing the swarm details
      swarmdetails , err := client.InspectService(msg.Actor.Attributes["com.docker.swarm.service.id"])
            if err != nil {
                log.Fatal(err)
            }
      fmt.Println(msg)

*/
      // end of testing


      //if(action == "start"){
      swarmdetails , err := client.InspectService(msg.Actor.Attributes["com.docker.swarm.service.id"])
      if(status == "container"){
        if(action == "start") {
          // lakukan checking service 
          if((getServiceCount(client,filter) != serviceCount) && (swarmdetails.Spec.Labels["nginxy.domain"] != "")){ // jika ada service baru, maka execute 


            if err != nil {
                log.Fatal(err)
            }

 /*          fmt.Println("---- name and labels goes down here -- ")
           fmt.Println("nama service : ",swarmdetails.Spec.Name)
           fmt.Println("hasil label : ",swarmdetails.Spec.Labels)
           fmt.Println("jumlah service : ",getServiceCount(client))

           //testing label swarm
           fmt.Println("label domain : ",swarmdetails.Spec.Labels["nginxy.uri"])

*/
           //-- define every var to make variable shorter

           name := swarmdetails.Spec.Name
           domain := swarmdetails.Spec.Labels["nginxy.domain"]
           port := swarmdetails.Spec.Labels["nginxy.port"]
           ssl := swarmdetails.Spec.Labels["nginxy.ssl"]
           key := swarmdetails.Spec.Labels["nginxy.ssl.key"]
           cert := swarmdetails.Spec.Labels["nginxy.ssl.cert"]

           //-- check the domain label is exist

           if (domain != ""){
             //-- assign to Nginx Struct to write the file configurations
             nginxconf := Nginx{ServiceName:name,DomainName:domain,ServicePort:port,Ssl:ssl,SslKey:key,SslCert:cert}
             fmt.Println(nginxconf)
             nginxconf.WriteConf() // write the nginx configuration based on the name
             reloadNginx()
           }else {
             log.Println("domain not set, please set")
           }

           // increment the service count
           serviceCount++

         }else{ // tak ada service baru
           log.Println("No New Services")
         }

        }else if(action == "destroy"){
          // -- destroy -> hapus container
          if(getServiceCount(client,filter) != serviceCount){ // jika jumlah total service berbeda dengan sebelumnya, maka terjadi penghapusan service
//            fmt.Println(msg)
            nginxconf := Nginx{ServiceName:msg.Actor.Attributes["com.docker.swarm.service.name"]}
            log.Println("penghapusan Service",nginxconf.ServiceName)
            nginxconf.DeleteConf() // delete the existing configuration based on the name
            reloadNginx()
            serviceCount-- // decrement the service count
//            fmt.Println((msg.Actor.Attributes["com.docker.swarm.service.name"]))
          }
        }
      }

/*
//        log.Println(msg.Action)
//        log.Println(msg.Type)
        log.Println(msg.Actor)
        fmt.Println("ini attribut nya : ",msg.Actor.Attributes)
        fmt.Println("testing : ", msg.Actor.Attributes["com.docker.swarm.service.id"])


        fmt.Println("ini coba spec element nya")
        //fmt.Println(client.InspectService(msg.Actor.Attributes["com.docker.swarm.service.id"]))
//        testing , err := client.InspectService(msg.Actor.Attributes["com.docker.swarm.service.id"])

        if err != nil {
            log.Fatal(err)
        }


       fmt.Println("---- name and labels goes down here -- ")
       fmt.Println(testing.Spec.Name)
       fmt.Println(testing.Spec.Labels)
       */

//    case <-timeout:
//        return
    }
}

}


// to get current service list
func getServiceCount(client *docker.Client,filter map[string][]string) int{

  localCount := 0
  imgs, err := client.ListServices(docker.ListServicesOptions{Filters: filter})
    if err != nil {
        panic(err)
    }
    for _,_ = range imgs {
      localCount++
    }
 return localCount
}
