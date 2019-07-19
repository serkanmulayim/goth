package certificateservice

// const (
// 	certificatesPath    = "/certificates"
// 	TLSCertificatesPath = certificatesPath + "/tls"
//
// 	defaultTLSCertBundlePath = certificatesPath + "/defaults/tls/cert"
//
// 	signingCertPathFormat    = certificatesPath + "/signing" + "/%v" ///alias
//
// 	initialTLSCertFSPath       = "./resources/init.crt"
// 	initialTLSPrivateKeyFSPAth = "./resources/init.key"
//
// 	defaultTLSCertFSPath       = "./resources/default.crt"
// 	defaultTLSPrivateKeyFSPath = "./resources/default.key"
//
// 	initialDefaultTLSPrivateKeyPassword = "password"
// )
//
// //KEEP only signing certificates
//
// //CheckDefaultTLSCertInBootstrap : checks if default tls cert
// func CheckDefaultTLSCertInBootstrap(etcdClient *clientv3.Client) string, string, error {
//
// 	ctx, cancel := context.WithTimeout(context.Background(), etcdclientservice.EtcdClientClientTimeout*time.Millisecond)
// 	certGR, err := etcdClient.Get(ctx, defaultTLSCertBundlePath)
// 	cancel()
// 	if err != nil {
// 		log.Fatalf("Error in getting default TLS cert from etcd %v:", err)
// 		return "", "", err // no need
// 	}
//
// 	if len(certGR.Kvs) == 0 {
// 		//if there is no default cert get it from disk
// 		pk, cert, err := GetSystemDefaultTLSCert(etcdClient)
//     x509.
//
// 		if err != nil {
// 			log.Fatalf("Error in loading the default TLS keypair. It does not exist in FS nor in etcd : %v. Exiting", err)
// 		}
// 		return pk, cert, nil
// 	}
//
// 	//Read File
// 	return nil
// }
//
// //GetSystemDefaultTLSCert returns private key and cert from resources/server.pfx. And stores it to etcd
// func GetSystemDefaultTLSCert(etcdClient *clientv3.Client) (error) {
// 	cerStr, err := ioutil.ReadFile(initialTLSCertFSPath)
//
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
//   pkStr, err := ioutil.ReadFile(initialTLSPrivateKeyFSPath)
//   if err != nil {
//     return nil, nil, err
//   }
//
// 	//write to etcd
//   etcdclientservice.EtcdPut(key, val, cli)
// 	return pkcs12.Decode(dat, initialDefaultTLSPrivateKeyPassword)
//
// }
