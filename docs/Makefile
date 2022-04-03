default:
	mdbook serve

build:
	-rm -rf book docs
	mdbook build
	-rm book/.gitignore
	-rm -rf book/.git

	mv book docs

clean: