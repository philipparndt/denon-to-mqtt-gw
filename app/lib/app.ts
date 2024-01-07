import cron from "node-cron"
import { getAppConfig } from "./config/config"
import { log } from "./logger"

import { connectMqtt, publish } from "./mqtt/mqtt-client"
// @ts-ignore
import * as Denon from "denon-client"

let eventSource: EventSource
let denonClient: any

type StateType = {
    power: string
    volume: number
}

const state: StateType = {
    power: "OFF",
    volume: 0
}

export const triggerFullUpdate = async () => {
    log.debug("triggerFullUpdate")
    state.volume = await denonClient.getVolume()
    state.power = await denonClient.getPower()
    publishState()
}

const publishState = () => {
    publish(state, "state")
}

const start = async () => {
    const ip = getAppConfig().denon.ip
    log.info("Connecting to Denon device on", ip)
    denonClient = new Denon.DenonClient(ip)
    await denonClient.connect()

    denonClient.on("masterVolumeChanged", (volume: any) => {
        state.volume = volume
        log.debug("masterVolumeChanged", volume)
        publishState()
    })

    denonClient.on("powerChanged", (power: any) => {
        state.power = power
        log.debug("powerChanged", power)
        publishState()
    })

    denonClient.on("error", async (error: any) => {
        log.error("error", error)
        state.power = "ERROR"
        publishState()

        await start()
    })

    await triggerFullUpdate()
}

export const startApp = async () => {
    const mqttCleanUp = await connectMqtt()
    await start()
    await triggerFullUpdate()
    log.info("Application is now ready.")

    log.info("Scheduling refresh.")
    const task = cron.schedule("0 0 * * *", triggerFullUpdate)
    task.start()

    return () => {
        mqttCleanUp()
        eventSource?.close()
        task.stop()
    }
}
