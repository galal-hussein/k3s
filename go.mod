module github.com/rancher/k3s

go 1.13

replace (
	github.com/Microsoft/hcsshim => github.com/Microsoft/hcsshim v0.8.9
	github.com/benmoss/go-powershell => github.com/rancher/go-powershell v0.0.0-20200701184732-233247d45373
	github.com/containerd/btrfs => github.com/containerd/btrfs v0.0.0-20181101203652-af5082808c83
	github.com/containerd/cgroups => github.com/containerd/cgroups v0.0.0-20200531161412-0dbf7f05ba59
	github.com/containerd/console => github.com/containerd/console v0.0.0-20181022165439-0650fd9eeb50
	github.com/containerd/containerd => github.com/rancher/containerd v1.4.0-k3s1
	github.com/containerd/continuity => github.com/containerd/continuity v0.0.0-20190815185530-f2a389ac0a02
	github.com/containerd/fifo => github.com/containerd/fifo v0.0.0-20190816180239-bda0ff6ed73c
	github.com/containerd/go-runc => github.com/containerd/go-runc v0.0.0-20200220073739-7016d3ce2328
	github.com/containerd/typeurl => github.com/containerd/typeurl v0.0.0-20180627222232-a93fcdb778cd
	github.com/coreos/flannel => github.com/rancher/flannel v0.12.0-k3s1
	github.com/coreos/go-systemd => github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20190205005809-0d3efadf0154
	github.com/docker/docker => github.com/docker/docker v17.12.0-ce-rc1.0.20190219214528-cbe11bdc6da8+incompatible
	github.com/docker/libnetwork => github.com/docker/libnetwork v0.8.0-dev.2.0.20190624125649-f0e46a78ea34
	github.com/golang/protobuf => github.com/golang/protobuf v1.3.5
	github.com/juju/errors => github.com/rancher/nocode v0.0.0-20200630202308-cb097102c09f
	github.com/kubernetes-sigs/cri-tools => github.com/rancher/cri-tools v1.19.0-k3s1
	github.com/matryer/moq => github.com/rancher/moq v0.0.0-20190404221404-ee5226d43009
	github.com/opencontainers/runc => github.com/opencontainers/runc v1.0.0-rc92
	github.com/opencontainers/runtime-spec => github.com/opencontainers/runtime-spec v1.0.3-0.20200728170252-4d89ac9fbff6
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20200224152610-e50cd9704f63
	google.golang.org/grpc => google.golang.org/grpc v1.27.1
	gopkg.in/square/go-jose.v2 => gopkg.in/square/go-jose.v2 v2.2.2
	k8s.io/api => github.com/rancher/kubernetes/staging/src/k8s.io/api v1.19.0-k3s1
	k8s.io/apiextensions-apiserver => github.com/rancher/kubernetes/staging/src/k8s.io/apiextensions-apiserver v1.19.0-k3s1
	k8s.io/apimachinery => github.com/rancher/kubernetes/staging/src/k8s.io/apimachinery v1.19.0-k3s1
	k8s.io/apiserver => github.com/rancher/kubernetes/staging/src/k8s.io/apiserver v1.19.0-k3s1
	k8s.io/cli-runtime => github.com/rancher/kubernetes/staging/src/k8s.io/cli-runtime v1.19.0-k3s1
	k8s.io/client-go => github.com/rancher/kubernetes/staging/src/k8s.io/client-go v1.19.0-k3s1
	k8s.io/cloud-provider => github.com/rancher/kubernetes/staging/src/k8s.io/cloud-provider v1.19.0-k3s1
	k8s.io/cluster-bootstrap => github.com/rancher/kubernetes/staging/src/k8s.io/cluster-bootstrap v1.19.0-k3s1
	k8s.io/code-generator => github.com/rancher/kubernetes/staging/src/k8s.io/code-generator v1.19.0-k3s1
	k8s.io/component-base => github.com/rancher/kubernetes/staging/src/k8s.io/component-base v1.19.0-k3s1
	k8s.io/cri-api => github.com/rancher/kubernetes/staging/src/k8s.io/cri-api v1.19.0-k3s1
	k8s.io/csi-translation-lib => github.com/rancher/kubernetes/staging/src/k8s.io/csi-translation-lib v1.19.0-k3s1
	k8s.io/kube-aggregator => github.com/rancher/kubernetes/staging/src/k8s.io/kube-aggregator v1.19.0-k3s1
	k8s.io/kube-controller-manager => github.com/rancher/kubernetes/staging/src/k8s.io/kube-controller-manager v1.19.0-k3s1
	k8s.io/kube-proxy => github.com/rancher/kubernetes/staging/src/k8s.io/kube-proxy v1.19.0-k3s1
	k8s.io/kube-scheduler => github.com/rancher/kubernetes/staging/src/k8s.io/kube-scheduler v1.19.0-k3s1
	k8s.io/kubectl => github.com/rancher/kubernetes/staging/src/k8s.io/kubectl v1.19.0-k3s1
	k8s.io/kubelet => github.com/rancher/kubernetes/staging/src/k8s.io/kubelet v1.19.0-k3s1
	k8s.io/kubernetes => github.com/rancher/kubernetes v1.19.0-k3s1
	k8s.io/legacy-cloud-providers => github.com/rancher/kubernetes/staging/src/k8s.io/legacy-cloud-providers v1.19.0-k3s1
	k8s.io/metrics => github.com/rancher/kubernetes/staging/src/k8s.io/metrics v1.19.0-k3s1
	k8s.io/node-api => github.com/rancher/kubernetes/staging/src/k8s.io/node-api v1.19.0-k3s1
	k8s.io/sample-apiserver => github.com/rancher/kubernetes/staging/src/k8s.io/sample-apiserver v1.19.0-k3s1
	k8s.io/sample-cli-plugin => github.com/rancher/kubernetes/staging/src/k8s.io/sample-cli-plugin v1.19.0-k3s1
	k8s.io/sample-controller => github.com/rancher/kubernetes/staging/src/k8s.io/sample-controller v1.19.0-k3s1
	mvdan.cc/unparam => mvdan.cc/unparam v0.0.0-20190209190245-fbb59629db34
)

require (
	bitbucket.org/bertimus9/systemstat v0.0.0-20180207000608-0eeff89b0690 // indirect
	github.com/JeffAshton/win_pdh v0.0.0-20161109143554-76bb4ee9f0ab // indirect
	github.com/Microsoft/go-winio v0.4.15-0.20190919025122-fc70bd9a86b5 // indirect
	github.com/NYTimes/gziphandler v1.1.1 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/alexflint/go-filemutex v0.0.0-20171022225611-72bdc8eae2ae // indirect
	github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e // indirect
	github.com/auth0/go-jwt-middleware v0.0.0-20170425171159-5493cabe49f7 // indirect
	github.com/benmoss/go-powershell v0.0.0-00010101000000-000000000000 // indirect
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/bronze1man/goStrongswanVici v0.0.0-20190828090544-27d02f80ba40 // indirect
	github.com/buger/jsonparser v0.0.0-20180808090653-f4dd9f5a6b44 // indirect
	github.com/clusterhq/flocker-go v0.0.0-20160920122132-2b8b7259d313 // indirect
	github.com/codegangsta/negroni v1.0.0 // indirect
	github.com/container-storage-interface/spec v1.2.0 // indirect
	github.com/containerd/aufs v0.0.0-20191030083217-371312c1e31c // indirect
	github.com/containerd/btrfs v0.0.0-20200117014249-153935315f4a // indirect
	github.com/containerd/cgroups v0.0.0-20200710171044-318312a37340 // indirect
	github.com/containerd/continuity v0.0.0-20200710164510-efbc4488d8fe // indirect
	github.com/containerd/cri v1.11.1-0.20200820101445-b0cc07999aa5
	github.com/containerd/fifo v0.0.0-20200410184934-f15a3290365b // indirect
	github.com/containerd/go-cni v1.0.0 // indirect
	github.com/containerd/go-runc v0.0.0-20200220073739-7016d3ce2328 // indirect
	github.com/containerd/ttrpc v1.0.1 // indirect
	github.com/containerd/typeurl v1.0.1 // indirect
	github.com/containerd/zfs v0.0.0-20191030014035-9abf673ca6ff // indirect
	github.com/containernetworking/cni v0.8.0 // indirect
	github.com/containernetworking/plugins v0.8.1 // indirect
	github.com/coredns/corefile-migration v1.0.10 // indirect
	github.com/coreos/go-iptables v0.4.2
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f
	github.com/d2g/dhcp4 v0.0.0-20170904100407-a1d1b6c41b1c // indirect
	github.com/d2g/dhcp4client v1.0.0 // indirect
	github.com/d2g/dhcp4server v0.0.0-20181031114812-7d4a0a7f59a5 // indirect
	github.com/d2g/hardwareaddr v0.0.0-20190221164911-e7d9fbe030e4 // indirect
	github.com/denverdino/aliyungo v0.0.0-20170629053852-f6cab0c35083 // indirect
	github.com/docker/docker v17.12.0-ce-rc1.0.20200821074627-7ae5222c72cc+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-events v0.0.0-20190806004212-e31b211e4f1c // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/erikdubbelboer/gspt v0.0.0-20190125194910-e68493906b83
	github.com/euank/go-kmsg-parser v2.0.0+incompatible // indirect
	github.com/frankban/quicktest v1.10.2 // indirect
	github.com/fullsailor/pkcs7 v0.0.0-20190404230743-d7302db945fa // indirect
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/go-ini/ini v1.28.1 // indirect
	github.com/go-ozzo/ozzo-validation v3.5.0+incompatible // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/godbus/dbus v0.0.0-20180201030542-885f9cc04c9c // indirect
	github.com/gogo/googleapis v1.3.2 // indirect
	github.com/google/cadvisor v0.35.0 // indirect
	github.com/google/tcpproxy v0.0.0-20180808230851-dfa16c61dad2
	github.com/google/uuid v1.1.1
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/gorilla/websocket v1.4.1
	github.com/hashicorp/go-multierror v1.0.0 // indirect
	github.com/heketi/heketi v9.0.1-0.20190917153846-c2e2a4ab7ab9+incompatible // indirect
	github.com/heketi/tests v0.0.0-20151005000721-f3775cbcefd6 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/ishidawataru/sctp v0.0.0-20190723014705-7c296d48a2b5 // indirect
	github.com/j-keck/arping v0.0.0-20160618110441-2cf9dc699c56 // indirect
	github.com/joho/godotenv v0.0.0-20161216230537-726cc8b906e3 // indirect
	github.com/juju/errors v0.0.0-20180806074554-22422dad46e1 // indirect
	github.com/juju/loggo v0.0.0-20190526231331-6e530bcce5d8 // indirect
	github.com/juju/testing v0.0.0-20190613124551-e81189438503 // indirect
	github.com/karrick/godirwalk v1.7.5 // indirect
	github.com/lib/pq v1.1.1
	github.com/libopenstorage/openstorage v1.0.0 // indirect
	github.com/lpabon/godbc v0.1.1 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mattn/go-shellwords v1.0.3 // indirect
	github.com/mattn/go-sqlite3 v1.13.0
	github.com/miekg/dns v1.1.4 // indirect
	github.com/mindprince/gonvml v0.0.0-20190828220739-9ebdce4bb989 // indirect
	github.com/mistifyio/go-zfs v2.1.2-0.20190413222219-f784269be439+incompatible // indirect
	github.com/moby/ipvs v1.0.1 // indirect
	github.com/mohae/deepcopy v0.0.0-20170603005431-491d3605edfb // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/mvdan/xurls v1.1.0 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/onsi/ginkgo v1.14.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v1.0.0-rc92
	github.com/opencontainers/runtime-tools v0.0.0-20181011054405-1d69bd0f9c39 // indirect
	github.com/opencontainers/selinux v1.6.0
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible
	github.com/pkg/errors v0.9.1
	github.com/quobyte/api v0.1.2 // indirect
	github.com/rakelkar/gonetsh v0.0.0-20190930180311-e5c5ffe4bdf0 // indirect
	github.com/rancher/dynamiclistener v0.2.1
	github.com/rancher/go-powershell v0.0.0-20200701182037-6845e6fcfa79 // indirect
	github.com/rancher/helm-controller v0.7.3
	github.com/rancher/kine v0.4.0
	github.com/rancher/remotedialer v0.2.0
	github.com/rancher/wrangler v0.6.1
	github.com/rancher/wrangler-api v0.6.0
	github.com/robfig/cron v1.1.0 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/rootless-containers/rootlesskit v0.10.0
	github.com/safchain/ethtool v0.0.0-20190326074333-42ed695e3de8 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/storageos/go-api v0.0.0-20180912212459-343b3eff91fc // indirect
	github.com/stretchr/testify v1.6.1
	github.com/tchap/go-patricia v2.3.0+incompatible // indirect
	github.com/thecodeteam/goscaleio v0.1.0 // indirect
	github.com/urfave/cli v1.22.2
	github.com/urfave/negroni v1.0.0 // indirect
	github.com/vishvananda/netns v0.0.0-20200520041808-52d707b772fe // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v0.0.0-20180618132009-1d523034197f // indirect
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200819165624-17cef6e3e9d5
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
	gonum.org/v1/gonum v0.6.2 // indirect
	gonum.org/v1/netlib v0.0.0-20190331212654-76723241ea4e // indirect
	google.golang.org/grpc v1.27.0
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce // indirect
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
	gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	k8s.io/api v0.19.0
	k8s.io/apimachinery v0.19.0
	k8s.io/apiserver v0.0.0
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/cloud-provider v0.0.0
	k8s.io/cluster-bootstrap v0.0.0 // indirect
	k8s.io/component-base v0.19.0
	k8s.io/cri-api v0.19.0
	k8s.io/heapster v1.2.0-beta.1 // indirect
	k8s.io/klog v1.0.0
	k8s.io/kube-controller-manager v0.0.0 // indirect
	k8s.io/kube-proxy v0.0.0 // indirect
	k8s.io/kube-scheduler v0.0.0 // indirect
	k8s.io/kubectl v0.0.0 // indirect
	k8s.io/kubelet v0.0.0 // indirect
	k8s.io/legacy-cloud-providers v0.0.0 // indirect
	k8s.io/sample-apiserver v0.0.0 // indirect
	k8s.io/system-validators v1.1.2 // indirect
	sigs.k8s.io/yaml v1.2.0
)
