all: dist

build:
	go build -o moonshot
dist: build
	mkdir moonshot_game
	cp moonshot moonshot_game/
	cp -r assets moonshot_game/
	tar cvfz moonshot.tar.gz moonshot_game
	rm -rf moonshot_game
	
