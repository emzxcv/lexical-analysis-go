
.PHONY: default setup-developer-tools

default: levenshtein soundex nysiis phonex

levenshtein:
	$(info [+] ===============Levenshtein===============)
	go run cmd/levenshtein.go

soundex:
	$(info [+] ===============Levenshtein-Soundex===========)
	go run cmd/soundex.go

nysiis:
	$(info [+] ===============Levenshtein-NYSIIS===========)
	go run cmd/nysiis.go

phonex:
	$(info [+] ===============Levenshtein-Phonex===========)
	go run cmd/phonex.go
