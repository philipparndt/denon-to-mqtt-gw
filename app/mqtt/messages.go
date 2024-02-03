package mqtt

type DeviceDataMessage struct {
    Interface   string  `json:"interface"`
    Detection   string  `json:"detection"`
    Status      int     `json:"status"`
    Class       string  `json:"class"`
    Energy      float64 `json:"energy"`
    Voltage     float64 `json:"voltage"`
    CurrentmA   float64 `json:"current_mA"`
    TotalWhr    float64 `json:"total_Whr"`
    Temperature float64     `json:"temperature"`
}

type TransmitMessage struct {
    Port       string `json:"port"`
    BytesTx    int64  `json:"bytesTx"`
    BytesRx    int64  `json:"bytesRx"`
    PacketsTx  int64  `json:"packetsTx"`
    PacketsRx  int64  `json:"packetsRx"`
    TotalBytes int64  `json:"totalBytes"`
}

type AggregatedEnergy struct {
    EnergySum float64 `json:"energySum"`
    WhrSum    float64 `json:"whrSum"`
}

