# mdoc

This is a simple Markdown server for exploring technical documentation.

The main and only product objective is to enable rich exploration of technical documentation (a la Google's g3doc).

## Work in progress

This is a work in progress. It's currently used inside Prey to view through the `doc/` directories within the global codebase.

The ultimate goal is to support the following features (in order of importance):

1. Search: smart, real-time search
2. Fast: should implement HTTP caches, hot pages, etc.
3. Secure: the server should only serve markdown pages within the defined base directory
4. Linkable: code references should be referenced to the corresponding repository
5. Quality: a user feedback system so that authors can review the perceived quality of their docs
6. Git support
7. Customizable: users should be able to edit the main HTML template
8. Interoperable: users should be able to use similar applications, such as GitBook, GitHub online editing, etc