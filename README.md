# PrivatePredict 🔒
A privacy-first prediction market built on Aleo.

## What is PrivatePredict?
PrivatePredict is an adaptation of [SocialPredict](https://github.com/openpredictionmarkets/socialpredict) — an open source prediction market engine — with a core privacy layer added using Aleo's zero-knowledge proof system.

Users place bets on prediction markets **privately**. Their outcome and amount are never visible on-chain. Only a ZK proof exists on the Aleo network.

## Privacy Guarantees
- ✅ Bet outcome is private — nobody knows if you bet YES or NO
- ✅ Bet amount is private — stored in encrypted BetRecord
- ✅ No address leaks — zero finalize blocks, pure records only
- ✅ credits.aleo integrated for on-chain payments

## Tech Stack
- Frontend: React
- Backend: Go
- Database: PostgreSQL
- ZK Layer: Aleo / Leo (socialpredict.aleo)
- Prover: Node.js + @provablehq/sdk

## Based On
- [SocialPredict](https://github.com/openpredictionmarkets/socialpredict) — MIT License
- [Aleo](https://aleo.org) — ZK blockchain
