module github.com/crossplane-contrib/provider-github

go 1.13

require (
	github.com/crossplane/crossplane-runtime v0.15.0
	github.com/crossplane/crossplane-tools v0.0.0-20201201125637-9ddc70edfd0d
	github.com/fatih/color v1.9.0 // indirect
	github.com/google/go-cmp v0.5.5
	github.com/google/go-github/v33 v33.0.0
	github.com/google/uuid v1.1.4 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mitchellh/copystructure v1.2.0
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	sigs.k8s.io/controller-runtime v0.9.2
	sigs.k8s.io/controller-tools v0.2.4
)
