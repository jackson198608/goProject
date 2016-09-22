package main

import (
		"fmt"
		"gopkg.in/olivere/elastic.v3"
		)


func main(){
	fmt.Println("begin connect to client\n");
	// Create a client
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.86.88:9200"));
	if err != nil {
    	// Handle error
	}
	fmt.Println(err);
	fmt.Println(client);

}
