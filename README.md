# GAE Git Renderer

Render a Git repo in the Google App Engine with some special tricks.

## Special files

Within each folder, a `.gaegr` file can be defined to describe what template to use to render the folder.
Templates that ship by default are:

- static - renders the folder as a static listing of files
- index - renders a index.html file in this folder


