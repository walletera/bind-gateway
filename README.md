# Walletera BIND Gateway
[![Go](https://github.com/walletera/bind-gateway/actions/workflows/go.yml/badge.svg)](https://github.com/walletera/bind-gateway/actions/workflows/go.yml)

This repository is part of the [Walletera project](https://github.com/walletera) and serves as the gateway that integrates the core Walletera platform with the API of [Banco Industrial](https://apibank.bind.com.ar/). Through this gateway, Walletera customers can send and receive payments using Argentina's [CVU](https://www.bcra.gob.ar/MediosPago/Politica_Pagos-i.asp#:~:text=Back%20to%20top-,Single%20Virtual%20Code%20(CVU),-What%20is%20a) (Clave Virtual Uniforme) immediate payments system.

## Overview

- **Purpose:**  
  To facilitate seamless integration between the Walletera financial platform and Banco Industrial's payment services, enabling CVU-based instant payments for users.

- **Features:**
    - Send payments to any CVU via Banco Industrial's API
    - Receive payments from any CVU into the Walletera ecosystem
    - Secure and reliable API communication
    - Real-time transaction updates and reconciliation

## Getting Started

1. **Clone the Repository**
```shell script
git clone https://github.com/walletera/walletera-banco-industrial-gateway.git
   cd walletera-banco-industrial-gateway
```

3. **Run the tests**

> **Note:**
You need docker running on your machine to be able to run the tests

```shell script
make test_all
```

## Contributing

We welcome contributions! Please open an issue or submit a pull request following our [contribution guidelines](CONTRIBUTING.md).

---

*For support, questions, or partnership opportunities, please contact the Walletera team.*