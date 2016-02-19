#ifndef LOGEVENT_H
#define LOGEVENT_H

#include <QEvent>

struct LogEvent : public QEvent
{
    enum {EventID = QEvent::User};
    LogEvent(const QString message_) : QEvent(static_cast<Type>(EventID)), message(message_) {}

    const QString message;
};

#endif // LOGEVENT_H
