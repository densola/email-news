emneBin = ./bin/email-news
ansibleFilesPath = ./ansible/roles/emne/files/
ansibleDir = ./ansible/

tgen:
	templ generate

run: tgen
	go run .

build: tgen
	go build -o $(emneBin)

deploy: build
	cp $(emneBin) $(ansibleFilesPath)
	ANSIBLE_CONFIG=$(ansibleDir) ansible-playbook $(ansibleDir)emne.yml --ask-become-pass