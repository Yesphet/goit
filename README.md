### Intro

Goit is a simple tool for writing standard git commit message by terminal user interface, and generate change log by commit log with one button. 

Format commit message with [AngularJs Git Commit Message Conventions](https://docs.google.com/document/d/1QrDFcIiPjSLDn3EL15IJygNPiHORgU1_OOAqWjiDU5Y/edit#heading=h.uyo6cb12dt6w) 

### How to use

#### 1. Installation

You can either download the binary directly from the downloads page or `go get` it:

```
go get -u github.com/Yesphet/goit
```

#### 2. Configuration

You can add `goit.yml` to any git repository.

It's an example config: 

```
---
commit:
    # Set scopes here, and it will autocomplete by these scopes when you do commit. 
    scopes: ["commit"]
    # Specify your custom change type.  
    # Use format: - "$name: $description"
    types: 
        - feat: new feature
        - bug: bug fix
``` 

#### 2. Set up git alias

Just run `git config --global alias.cz '!goit cz'`, then `cd` into any git repository and use `git cz` instead of `git commit` and you will find the commitizen prompt. 

