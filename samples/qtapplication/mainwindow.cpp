#include "mainwindow.h"
#include "ui_mainwindow.h"
#include <QLocalSocket>

MainWindow::MainWindow(QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    this->server = nullptr;
}

MainWindow::~MainWindow()
{
    delete ui;
}

void MainWindow::on_LaunchClient_clicked()
{
    this->clientSocket = new QLocalSocket(this);
    connect(this->clientSocket, SIGNAL(disconnected()), this->clientSocket, SLOT(deleteLater()));
    connect(this->clientSocket, SIGNAL(readyRead()), this, SLOT(onReadFromServer()));
    connect(this->clientSocket, SIGNAL(error(QLocalSocket::LocalSocketError)),
                        this, SLOT(onSocketError(QLocalSocket::LocalSocketError)));
    this->clientSocket->connectToServer("testlocalsocket");

    auto sendCount = qrand() % 20;
    for (int i = 0; i < sendCount; i++) {
        this->ui->Log->appendPlainText(QString("Writing '%1'").arg(i));
        this->clientSocket->write((QString::number(i) + "\n").toUtf8());
    }
    this->clientSocket->write("end\n");
    this->clientSocket->flush();
}

void MainWindow::on_RunServer_clicked()
{
    if (this->server != nullptr) {
        this->ui->Log->appendPlainText("Server is running");
        return;
    }
    this->server = new QLocalServer(this);
    this->ui->Log->appendPlainText("Starting server\n");
    this->server->listen("testlocalsocket");
    connect(this->server, SIGNAL(newConnection()), this, SLOT(onConnected()));
}

void MainWindow::on_StopServer_clicked()
{
    if (this->server == nullptr) {
        this->ui->Log->appendPlainText("No server is running");
        return;
    }
    this->server->close();
    this->ui->Log->appendPlainText("Stop server\n");
    delete this->server;
    this->server = nullptr;
}

void MainWindow::onConnected()
{
    this->ui->Log->appendPlainText("Connected from Client\n");
    this->serverSocket = this->server->nextPendingConnection();
    connect(this->serverSocket, SIGNAL(disconnected()), this->serverSocket, SLOT(deleteLater()));
    connect(this->serverSocket, SIGNAL(readyRead()), this, SLOT(onReadFromClient()));
    connect(this->serverSocket, SIGNAL(error(QLocalSocket::LocalSocketError)),
                        this, SLOT(onSocketError(QLocalSocket::LocalSocketError)));

    auto sendCount = qrand() % 20;
    for (int i = 0; i < sendCount; i++) {
        this->ui->Log->appendPlainText(QString("Writing '%1'").arg(i));
        this->serverSocket->write((QString::number(i) + "\n").toUtf8());
    }
    this->serverSocket->write("end\n");
    this->serverSocket->flush();
}

void MainWindow::onReadFromClient()
{
    for (int i = 0; i < 100; i++) {
        if (!this->serverSocket->canReadLine()) {
            break;
        }
        auto buffer = this->serverSocket->readLine();
        auto received = QString(buffer).trimmed();
        this->ui->Log->appendPlainText(QString("Read '%1'").arg(received));
        if (received == "end") {
            this->serverSocket->disconnect();
            this->serverSocket = nullptr;
            this->ui->Log->appendPlainText("disconnected");
            break;
        }
    }
}

void MainWindow::onReadFromServer()
{
    for (int i = 0; i < 100; i++) {
        if (!this->clientSocket->canReadLine()) {
            break;
        }
        auto buffer = this->clientSocket->readLine();
        auto received = QString(buffer).trimmed();
        this->ui->Log->appendPlainText(QString("Read '%1'").arg(received));
        if (received == "end") {
            this->clientSocket->disconnect();
            this->clientSocket = nullptr;
            this->ui->Log->appendPlainText("disconnected");
            break;
        }
    }
}

void MainWindow::onSocketError(QLocalSocket::LocalSocketError socketError)
{
    switch (socketError) {
    case QLocalSocket::ServerNotFoundError:
        this->ui->Log->appendPlainText("Server Not Found\n");
        break;
    case QLocalSocket::ConnectionRefusedError:
        this->ui->Log->appendPlainText("Connection Refused\n");
        break;
    case QLocalSocket::PeerClosedError:
        break;
    default:
        this->ui->Log->appendPlainText(QString("Error: %1").arg("Connection Refused\n"));
    }
}


void MainWindow::on_MainWindow_destroyed()
{
    if (this->server != nullptr) {
        this->server->close();
    }
}
