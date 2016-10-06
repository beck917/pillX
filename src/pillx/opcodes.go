package pillx

const SYS_CONNECT_MARK_ERROR uint16 = 0X0001
const SYS_CONNECT_SIZE_ERROR uint16 = 0X0002
const SYS_CONNECT_HANDSHAKE_ERROR uint16 = 0X0003
const SYS_CONNECT_WORKER_ERROR uint16 = 0X0004

const SYS_ON_CONNECT uint16 = 0X0010
const SYS_ON_MESSAGE uint16 = 0X0011
const SYS_ON_CLOSE uint16 = 0X0012
const SYS_ON_HANDSHAKE uint16 = 0X0013
const SYS_ON_BLOCK uint16 = 0X0014     //屏蔽
const SYS_ON_BAN uint16 = 0X0015       //禁言
const SYS_ON_KICK uint16 = 0X0016      //踢
const SYS_ON_BLACK uint16 = 0X0017     //黑名单
const SYS_ON_HEARTBEAT uint16 = 0X0018 //心跳
const SYS_ON_CLIENTIN uint16 = 0X0019  //广播进入
const SYS_ON_CLIENTOUT uint16 = 0X0020 //广播退出

const PROTO_HEADER_FIRSTCHAR = 0x7f
