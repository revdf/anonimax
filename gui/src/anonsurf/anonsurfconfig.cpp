#include "anonsurf/anonsurfconfig.h"

#include <QDateTime>
#include <QDir>
#include <QFile>
#include <QProcess>
#include <QStandardPaths>
#include <QTextStream>

namespace {
constexpr auto AnonSurfBin = "/usr/bin/anonimax";
} // namespace

AnonSurfConfig::AnonSurfConfig(QObject *parent) : QObject(parent) {
  connect(&m_writeProcess,
          QOverload<int, QProcess::ExitStatus>::of(&QProcess::finished), this,
          [this](int, QProcess::ExitStatus) {
            m_busy = false;
            refresh();
          });
  refresh();
}

void AnonSurfConfig::refresh() {
  QProcess proc;
  proc.start(AnonSurfBin, {QStringLiteral("config"), QStringLiteral("list")});
  proc.waitForFinished(3000);
  if (proc.exitCode() == 0)
    parseOutput(QString::fromUtf8(proc.readAllStandardOutput()));
  emit changed();
}

void AnonSurfConfig::setSetting(const QString &key, const QString &value) {
  if (m_busy)
    return;
  m_busy = true;
  emit changed();
  m_writeProcess.start(
      QStringLiteral("pkexec"),
      {AnonSurfBin, QStringLiteral("config"), QStringLiteral("set"), key,
       value});
}

bool AnonSurfConfig::saveLog(const QString &text) {
  const QString dir = QStandardPaths::writableLocation(
                          QStandardPaths::AppLocalDataLocation)
                      + QStringLiteral("/logs");
  if (!QDir().mkpath(dir))
    return false;
  const QString filename = dir + QStringLiteral("/anonimax-")
                           + QDateTime::currentDateTime().toString(
                                 QStringLiteral("yyyy-MM-dd_HH-mm-ss"))
                           + QStringLiteral(".log");
  QFile f(filename);
  if (!f.open(QIODevice::WriteOnly | QIODevice::Text))
    return false;
  QTextStream out(&f);
  out << text;
  return true;
}

void AnonSurfConfig::parseOutput(const QString &output) {
  for (const QString &line :
       output.split(QLatin1Char('\n'), Qt::SkipEmptyParts)) {
    const int eq = line.indexOf(QLatin1Char('='));
    if (eq < 0)
      continue;
    const QString k = line.left(eq).trimmed();
    const QString v = line.mid(eq + 1).trimmed();
    const bool bv   = (v == QLatin1String("true"));

    if      (k == QLatin1String("change_mac"))           m_changeMac        = bv;
    else if (k == QLatin1String("renew_ip"))             m_renewIP          = bv;
    else if (k == QLatin1String("force_utc"))            m_forceUTC         = bv;
    else if (k == QLatin1String("block_webrtc"))         m_blockWebRTC      = bv;
    else if (k == QLatin1String("block_telemetry"))      m_blockTelemetry   = bv;
    else if (k == QLatin1String("disable_camera_mic"))   m_disableCameraMic = bv;
    else if (k == QLatin1String("clean_on_stop"))        m_cleanOnStop      = bv;
    else if (k == QLatin1String("use_vpn"))              m_useVPN           = bv;
    else if (k == QLatin1String("use_proxy"))            m_useProxy         = bv;
    else if (k == QLatin1String("vpn_config_path"))      m_vpnConfigPath    = v;
    else if (k == QLatin1String("socks5_proxy_address")) m_socks5Address    = v;
    else if (k == QLatin1String("http_proxy_address"))   m_httpAddress      = v;
  }
}
