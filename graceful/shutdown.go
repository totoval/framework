package graceful

func ShutDown(quietly bool) {
	logInfo(quietly, "Totoval is shutting down")
	closeQueue(quietly)
	closeCache(quietly)
	closeDB(quietly)
	closeMonitor(quietly)
	logInfo(quietly, "Totoval is shut down")
}
