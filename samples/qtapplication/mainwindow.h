#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <QLocalServer>
#include <QLocalSocket>
#include <QPlainTextEdit>

namespace Ui {
class MainWindow;
}

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    explicit MainWindow(QWidget *parent = 0);
    ~MainWindow();

private slots:
    void on_LaunchClient_clicked();
    void on_RunServer_clicked();
    void on_StopServer_clicked();
    void onConnected();
    void onReadFromClient();
    void onReadFromServer();
    void onSocketError(QLocalSocket::LocalSocketError socketError);

    void on_MainWindow_destroyed();

private:
    Ui::MainWindow *ui;
    QLocalServer *server;
    QLocalSocket *serverSocket;
    QLocalSocket *clientSocket;
};

#endif // MAINWINDOW_H
