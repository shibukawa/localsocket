LocalSocket
=================

This module provides a counter part of `QLocalServer <http://doc.qt.io/qt-5/qlocalserver.html>`_ and `QLocalSocket <http://doc.qt.io/qt-5/qlocalsocket.html>`_. It is designed for implementing core logic layer by Golang behind Qt GUI.

As described in Qt's document, it uses Unix domain socket on POSIX systems and Named Pipe on Windows.

Usage
==========

* ``localsocket.LocalSocket``

  It implements ``net.Conn`` and ``io.Reader`` and ``io.Writer`` interface.

  * ``Read(data []byte) (int, error)``
  * ``Write(data []byte) (int, error)``
  * ``Close() error``
  * ``LocalAddr() net.Addr``
  * ``RemoteAddr() net.Addr``
  * ``SetDeadline(t time.Time) error``
  * ``SetReadDeadline(t time.Time) error``
  * ``SetWriteDeadline(t time.Time) error``

* ``localsocket.NewLocalSocket(address string) (*LocalSocket, error)``

  It connects via Unix Domain Socket at ``${TMP}/<address>`` on POSIX, or Named Pipe at ``\\.\pipe\<address>`` on Windows .

* ``localsocket.NewLocalServer(address string) *LocalServer``

  * ``Listen() error``
    
    Creates Unix Domain Socket at ``${TMP}/<address>`` on POSIX, or Named Pipe at ``\\.\pipe\<address>`` on Windows and start waiting client.

  * ``ListenAndServe() error``

    It is similar to ``Listen()``, but blocks execution until ``Close()`` is called.

  * ``SetOnConnectionCallback(callback func(socket *LocalSocket))``

    Registers callback function that is called when new client connects. 

  * ``Close()``

    Closes connection.

  * ``Path() string``

    Returns actual path to Unix Domain Socket or Named Pipe.
