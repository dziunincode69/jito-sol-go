# Solana Raydium Token Sniper - JITO MEV Integration

This project integrates Solana's **Raydium Token Sniper** with **JITO MEV** for optimized transaction performance.  
However, **JITO recently changed their terms for mempool transactions**, which impacts the usage of this feature.  

### **Current Solution**  
You can **still use this tool** by **removing the JITO MEV function** and utilizing **public RPC providers** like:
- [Helius](https://helius.xyz/)
- [RPCPool](https://rpcpool.com/)

---

## Features
- Token sniping on Raydium pools with rapid transaction capabilities.
- Configurable RPC endpoint for quick customization.
- Supports custom wallet integration for smooth sniping transactions.

## Prerequisites
Ensure you have the following installed:
- **Go** programming language (version >= 1.18)
- **Solana CLI** configured with your wallet.
- **A valid RPC endpoint** (Helius, RPCPool, or others).

## Installation
Clone the repository:

```bash
git clone https://github.com/your-username/solana-jito-sniper-go.git
cd solana-jito-sniper-go
