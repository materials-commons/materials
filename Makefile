GODIRS = . wsmaterials materials site autoupdate send db
P = github.com/materials-commons/materials

all: fmt test
	(cd materials; go build materials.go)

test:
	rm -rf test_data/.materials
	rm -rf test_data/corrupted
	mkdir -p test_data/.materials/projects
	cp test_data/*.project test_data/.materials/projects
	cp test_data/.user test_data/.materials
	mkdir -p test_data/corrupted/.materials
	cp test_data/projects_corrupted test_data/corrupted/.materials/projects
	mkdir -p /tmp/tproj/a
	touch /tmp/tproj/a/a.txt
	-./runtests.sh
	rm -rf /tmp/tproj

fmt:
	-for d in $(GODIRS); do (cd $$d; go fmt); done
