package master

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	. "yu/common"
	. "yu/node"
	"yu/txn"
)

func getQryInfoFromReq(req *http.Request, params JsonString) (qcall *Qcall, err error) {
	tripodName, qryName := GetTripodCallName(req)
	blockNum, err := GetBlockNumber(req)
	if err != nil {
		return
	}
	qcall = &Qcall{
		TripodName:  tripodName,
		QueryName:   qryName,
		Params:      params,
		BlockNumber: blockNum,
	}
	return
}

func getExecInfoFromReq(req *http.Request, params JsonString) (tripodName, execName string, stxn txn.IsignedTxn, err error) {
	tripodName, callName := GetTripodCallName(req)
	ecall := &Ecall{
		TripodName: tripodName,
		ExecName:   callName,
		Params:     params,
	}
	caller := GetAddress(req)
	pubkey, sig, err := GetPubkeyAndSignature(req)
	if err != nil {
		return
	}
	stxn, err = txn.NewSignedTxn(caller, ecall, pubkey, sig)
	return
}

func getHttpJsonParams(c *gin.Context) (params JsonString, err error) {
	if c.Request.Method == http.MethodPost {
		params, err = readPostBody(c.Request.Body)
		if err != nil {
			return
		}
	} else {
		params = c.GetString(PARAMS_KEY)
	}
	return
}

func forwardQueryToWorker(ip string, rw http.ResponseWriter, req *http.Request) {
	director := func(req *http.Request) {
		req.URL.Host = ip
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(rw, req)
}
