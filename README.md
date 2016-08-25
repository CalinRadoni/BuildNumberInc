# BuildNumberInc

<p align="center">
__BuildNumberInc__ is an application to auto-increment the build number of a C/C++ project
</p>

### Usage

Add the executable to the pre-build step of your IDE.
Pass as parameter the full name of the source file that should be analyzed.

### Example

#### Atollic TrueSTUDIO for ARM

_project properties_ -> C/C++ Build -> Settings -> Build Steps -> Pre-build steps -> Command
```
<path_to_BuildNumberInc.exe>\BuildNumberInc.exe <path_to_source_file>\version.h
```

#### Version file

```cpp
#ifndef VERSION_H
#define VERSION_H

namespace YourNamespace {

#define SW_VER_MAJOR	1
#define SW_VER_MINOR	8
#define SW_VER_BUILD	28		// build number

} /* namespace */

#endif // VERSION_H
```
