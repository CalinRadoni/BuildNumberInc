# BuildNumberInc

__BuildNumberInc__ is an application to auto-increment the build number of a C/C++ project

## Note

This repository is archived because I have no use for this application any more.

## Usage

Add the executable and the arguments to the pre-build step of your IDE.

If the build number is declared as a _#define_ no flags are needed.
If the build number is declared as a _const_ use the __-c__ flag.
If you want to use only LF (and not CRLF) as line endings in the source file use __-r__ flag.
If you want to see the changed build number add the __-v__ flag.

### Example

#### Atollic TrueSTUDIO for ARM

In _project properties_ -> C/C++ Build -> Settings -> Build Steps -> Pre-build steps -> Command

for a _#define_ set:
```
<path_to_BuildNumberInc.exe>\BuildNumberInc.exe -v <path_to_file>\version.h SW_VER_BUILD
```
for a _const_ set:
```
<path_to_BuildNumberInc.exe>\BuildNumberInc.exe -c -v <path_to_file>\version.h verBuildNo_0
```

#### Test file (version.h)

```cpp
#ifndef VERSION_H
#define VERSION_H

#define SW_VER_MAJOR	1
#define SW_VER_MINOR	8
#define SW_VER_BUILD	28		// build number

namespace YourNamespace {

    const unsigned int verBuildNo_0 = 14; // a comment
    const uint16_t verBuildNo_1 = 27; // another comment

} // namespace

#endif // VERSION_H
```

## Building the package

__BuildNumberInc__ is a simple package so, if you have [Go](https://golang.org/) installed and
added to your path, just type `go build` in the source directory.

## License

__BuildNumberInc__ is released under the [MIT License](https://opensource.org/licenses/MIT).
