# yard-planning

Layanan backend untuk proses **suggestion**, **pickup**, dan **placement** container pada sistem penataan yard.  
Dibangun menggunakan **Golang + Echo**.

---

## 1. Requirements

Pastikan sudah terinstall:

### Golang
- Minimal versi **1.20+**
- Cek versi:
  ```bash
  go version

##2. Suggestion
Endpoint: /api/suggestion 
Digunakan untuk merencanakan planning penempatan cointainer
requestnya 

{
    "yard": "YRD1",
    "container_number": "ALFI000004",
    "container_size": 40,
    "container_height": 8.6,
    "container_type": "DRY"
}



##3 Placement
Endpoint: /api/placement
Digunakan untuk menempatkan container ketika setelah mendapatkan perencanaan penempatan
Requestnya
{
    "yard": "YRD1",
    "container_number": "ALFI000006",
    "block": "LC01",
    "slot": 1,
    "row": 1,
    "tier": 1
}


##4 Pickup
Enpoint: /api/pickup
Digunakan untuk mengeluarkan container dari tempat penyimpanannya
Requestnya 
{
    "yard": "YRD1",
    "container_number": "ALFI000003"
}
