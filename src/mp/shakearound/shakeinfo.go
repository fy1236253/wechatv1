package shakearound

import ()

/*
{
    "data": {
        "page_id ": 14211,
        "beacon_info": {
            "distance": 55.00620700469034,
            "major": 10001,
            "minor": 19007,
            "uuid": "FDA50693-A4E2-4FB1-AFCF-C6EB07647825"
        },
        "openid": "oVDmXjp7y8aG2AlBuRpMZTb1-cmA",
        " poi_id":1234
    },
    "errcode": 0,
    "errmsg": "success."
}
*/

type BeaconInfo struct {
	Distance float64 `json:"distance"` // Beacon信号与手机的距离，单位为米
	UUID     string  `json:"uuid"`
	Major    int     `json:"major"`
	Minor    int     `json:"minor"`
}

type Shakeinfo struct {
	PageId     int64      `json:"page_id"`     // 摇周边页面唯一ID
	BeaconInfo BeaconInfo `json:"beacon_info"` // 设备信息，包括UUID、major、minor，以及距离
	Openid     string     `json:"openid"`      // 商户AppID下用户的唯一标识
	PoiId      *int64     `json:"poi_id"`      // 门店ID，有的话则返回，反之不会在JSON格式内
}
