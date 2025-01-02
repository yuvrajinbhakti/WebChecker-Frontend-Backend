package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type DomainURL struct {
	DomainURL string `json:"domainurl"`
}

type DomainVar struct {
	Domain      string `json:"domain"`
	HasMX       bool   `json:"hasMX"`
	HasSPF      bool   `json:"hasSPF"`
	SPFRecord   string `json:"spfRecord"`
	HasDMARC    bool   `json:"hasDMARC"`
	DMARCRecord string `json:"dmarcRecord"`
}

var domainVars []DomainVar

func formHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    var domainUrl DomainURL
    err := json.NewDecoder(r.Body).Decode(&domainUrl)
    if err != nil {
        log.Printf("Error decoding request: %v\n", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    domainVar := isValidDomain(domainUrl.DomainURL)
    domainVars = append(domainVars, domainVar)
    json.NewEncoder(w).Encode(domainVars)
}

func isValidDomain(domain string) DomainVar {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	if len(mxRecords) > 0 {
		hasMX = true
	}
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}
	fmt.Printf("domain=%v\n,hasMX=%v\n,hasSPF=%v\n,spfRecord=%v\n,hasDMARC=%v\n,dmarcRecord=%v\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
	var domainVar DomainVar
	domainVar.Domain = domain
	domainVar.HasMX = hasMX
	domainVar.HasSPF = hasSPF
	domainVar.SPFRecord = spfRecord
	domainVar.HasDMARC = hasDMARC
	domainVar.DMARCRecord = dmarcRecord

	return domainVar
}

func main() {
    r := mux.NewRouter()
    
    // Add global CORS middleware
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
            
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    })
    
    r.HandleFunc("/form", formHandler).Methods("POST", "OPTIONS")
    fmt.Print("Server is running on port 8080\n")
    log.Fatal(http.ListenAndServe(":8080", r))
}
