emneBin = ./bin/email-news
emneDB = ./db.db
emneEnv = ./.env
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
	cp $(emneDB) $(ansibleFilesPath)
	cp $(emneEnv) $(ansibleFilesPath)
	ANSIBLE_CONFIG=$(ansibleDir) ansible-playbook $(ansibleDir)emne.yml --ask-become-pass