.PHONY: keygen

keygen:
	@KEY=$$(openssl ecparam -genkey -name secp256k1 -noout); \
	echo "Private Key: $$(echo "$$KEY" | openssl ec -text -noout 2>/dev/null | grep -A 3 'priv:' | tail -n +2 | tr -d '\n: ' | sed 's/[^0-9a-f]//g')"; \
	echo "Public Key: $$(echo "$$KEY" | openssl ec -text -noout 2>/dev/null | grep -A 5 'pub:' | tail -n +2 | tr -d '\n: ' | sed 's/[^0-9a-f]//g')"