warp is a dynamic analysis tool for Go. Or it might be, someday, if I finish it. For now, it is an **untested, incomplete, undocumented toy**.

Right now, warp helps detect mistaken assumptions about io.Reader.

To use, run `warp` over Go packages. **This will overwrite the package. Use source control.** Then run your unit tests.
