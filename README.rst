================
Go QCache-client
================

.. _QCache: https://github.com/tobgu/qcache

Go client library for QCache_. Uses consistent hashing to distribute data over multiple nodes.

Installation
============
::

    go get github.com/tobgu/go-qcache-client/qclient

Documentation
=============

Documentation is close to non-existent right now. Please see the tests in qclient_test.go for examples of how to use it.


Contributing
============
Want to contribute? That's great!

If you experience problems please log them on GitHub. If you want to contribute code,
please fork the code and submit a pull request.

If you intend to implement major features or make major changes please raise an issue
so that we can discuss it first.

TODO
====
- Timeouts and error handling
- Query method similar to that in the python client
- Documentation
- Possibility to POST large queries
