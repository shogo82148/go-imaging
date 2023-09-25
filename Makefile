.PHONY: examples
examples:
	go run ./internal/cmd/sample_gen examples

.PHONY: golden
golden:
	go run ./resize/internal/cmd/golden ./resize/testdata
	go run ./srgb/internal/cmd/golden ./srgb/testdata
	go run ./srgb/internal/cmd/decodeTone/main.go testdata/senkakuwan.png srgb/testdata/senkakuwan.golden.png testdata/gimp-linear.icc
