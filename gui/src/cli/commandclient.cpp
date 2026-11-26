#include "cli/commandclient.h"

#include <QDebug>
#include <QJsonArray>
#include <QJsonDocument>
#include <QJsonObject>

namespace {
constexpr auto AnonSurfPath = "/usr/bin/anonimax";
constexpr auto SafeKillPath = "/usr/lib/anonimax/safekill";
} // namespace

CommandClient::CommandClient(QObject *parent) : QObject(parent) {
  connect(&m_statusProcess, &QProcess::finished, this,
          &CommandClient::handleStatusFinished);
  connect(&m_ipProcess, &QProcess::finished, this,
          &CommandClient::handleIPFinished);
  connect(&m_safeKillProcess, &QProcess::finished, this,
          &CommandClient::handleSafeKillFinished);

  connect(&m_statusProcess, &QProcess::errorOccurred, this,
          [this](QProcess::ProcessError error) {
            if (error == QProcess::FailedToStart)
              emit statusFailed("The Anonimax CLI could not be started.");
          });
  connect(&m_ipProcess, &QProcess::errorOccurred, this,
          [this](QProcess::ProcessError error) {
            if (error == QProcess::FailedToStart)
              emit ipFailed("The Anonimax CLI could not be started.");
          });
  connect(&m_safeKillProcess, &QProcess::errorOccurred, this,
          [this](QProcess::ProcessError error) {
            if (error == QProcess::FailedToStart) {
              emit safeKillFinished(false, "SafeKill could not be started.");
            }
          });
  connect(&m_panicProcess, &QProcess::finished, this,
          &CommandClient::handlePanicFinished);
  connect(&m_panicProcess, &QProcess::errorOccurred, this,
          [this](QProcess::ProcessError error) {
            if (error == QProcess::FailedToStart)
              emit panicFinished(false, "Panic command could not be started.");
          });
  connect(&m_autochangeProcess, &QProcess::finished, this,
          &CommandClient::handleAutochangeFinished);
  connect(&m_autochangeProcess, &QProcess::errorOccurred, this,
          [this](QProcess::ProcessError error) {
            if (error == QProcess::FailedToStart)
              emit autochangeFinished(false, m_autochangeStarting,
                                      "Autochange command could not be started.");
          });
}

CommandClient::~CommandClient() {
  for (auto *process : {&m_statusProcess, &m_ipProcess, &m_safeKillProcess,
                        &m_panicProcess, &m_autochangeProcess}) {
    process->disconnect(this);
    if (process->state() != QProcess::NotRunning) {
      process->kill();
      if (!process->waitForFinished(1000))
        qWarning() << "Cannot finish GUI command shutdown:"
                   << process->errorString();
    }
  }
}

void CommandClient::requestStatus() {
  if (m_statusProcess.state() != QProcess::NotRunning)
    return;
  m_statusProcess.start(AnonSurfPath, {"status", "--json"});
}

void CommandClient::requestIP() {
  if (m_ipProcess.state() != QProcess::NotRunning)
    return;
  m_ipProcess.start(AnonSurfPath, {"ip", "--json"});
}

void CommandClient::runSafeKill() {
  if (m_safeKillProcess.state() != QProcess::NotRunning)
    return;
  m_safeKillProcess.start(SafeKillPath);
}

void CommandClient::handleStatusFinished(int exitCode,
                                         QProcess::ExitStatus exitStatus) {
  if (exitStatus != QProcess::NormalExit || exitCode != 0) {
    emit statusFailed(
        processError(m_statusProcess, "Anonimax status is unavailable."));
    return;
  }

  QJsonParseError parseError;
  const QJsonDocument document = QJsonDocument::fromJson(
      m_statusProcess.readAllStandardOutput(), &parseError);
  if (parseError.error != QJsonParseError::NoError || !document.isObject()) {
    emit statusFailed("Anonimax returned an invalid status response.");
    return;
  }

  const QJsonObject object = document.object();
  if (!object.value("anonimax_active").isBool() ||
      !object.value("tor_active").isBool() ||
      !object.value("dns_state").isString() ||
      !object.value("dns_addresses").isArray()) {
    emit statusFailed("Anonimax returned an incomplete status response.");
    return;
  }

  QStringList addresses;
  for (const QJsonValue &value : object.value("dns_addresses").toArray()) {
    if (value.isString())
      addresses.append(value.toString());
  }
  emit statusReady(object.value("anonimax_active").toBool(),
                   object.value("tor_active").toBool(),
                   object.value("dns_state").toString(), addresses);
}

void CommandClient::handleIPFinished(int exitCode,
                                     QProcess::ExitStatus exitStatus) {
  if (exitStatus != QProcess::NormalExit || exitCode != 0) {
    emit ipFailed(processError(m_ipProcess, "The Tor check failed."));
    return;
  }

  QJsonParseError parseError;
  const QJsonDocument document =
      QJsonDocument::fromJson(m_ipProcess.readAllStandardOutput(), &parseError);
  const QJsonObject object = document.object();
  if (parseError.error != QJsonParseError::NoError || !document.isObject() ||
      !object.value("ip").isString() || !object.value("is_tor").isBool()) {
    emit ipFailed("Anonimax returned an invalid Tor check response.");
    return;
  }
  emit ipReady(object.value("ip").toString(), object.value("is_tor").toBool());
}

void CommandClient::handleSafeKillFinished(int exitCode,
                                           QProcess::ExitStatus exitStatus) {
  const bool success = exitStatus == QProcess::NormalExit && exitCode == 0;
  emit safeKillFinished(
      success,
      success
          ? QString{}
          : processError(m_safeKillProcess,
                         "SafeKill failed. Ensure BleachBit is installed."));
}

void CommandClient::runPanic() {
  if (m_panicProcess.state() != QProcess::NotRunning)
    return;
  m_panicProcess.start(AnonSurfPath, {"panic", "--force"});
}

void CommandClient::runAutochange(bool start, int minutes) {
  if (m_autochangeProcess.state() != QProcess::NotRunning)
    return;
  m_autochangeStarting = start;
  if (start)
    m_autochangeProcess.start(AnonSurfPath,
                               {"autochange", QString::number(minutes)});
  else
    m_autochangeProcess.start(AnonSurfPath, {"autochange", "stop"});
}

void CommandClient::handlePanicFinished(int exitCode,
                                        QProcess::ExitStatus exitStatus) {
  const bool success = exitStatus == QProcess::NormalExit && exitCode == 0;
  emit panicFinished(
      success,
      success ? QString{} : processError(m_panicProcess, "Panic command failed."));
}

void CommandClient::handleAutochangeFinished(int exitCode,
                                             QProcess::ExitStatus exitStatus) {
  const bool success = exitStatus == QProcess::NormalExit && exitCode == 0;
  emit autochangeFinished(
      success, m_autochangeStarting,
      success ? QString{}
              : processError(m_autochangeProcess, "Autochange command failed."));
}

QString CommandClient::processError(QProcess &process,
                                    const QString &fallback) {
  const QString stderrText =
      QString::fromUtf8(process.readAllStandardError()).trimmed();
  if (stderrText.isEmpty())
    return fallback;

  const QStringList lines = stderrText.split('\n', Qt::SkipEmptyParts);
  for (auto line = lines.crbegin(); line != lines.crend(); ++line) {
    if (line->startsWith("[FAIL] "))
      return line->mid(7).trimmed();
  }
  return stderrText;
}
