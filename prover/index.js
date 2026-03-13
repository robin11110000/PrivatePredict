import express from 'express';
import { Account, ProgramManager, AleoNetworkClient, NetworkRecordProvider } from '@provablehq/sdk';

const app = express();
app.use(express.json());

const PROGRAM = `program PrivatePredict.aleo;

function place_bet:
    input r0 as u64.private;
    input r1 as bool.private;
    input r2 as u64.private;
    input r3 as address.public;
    output r0 as u64.private;
    output r1 as bool.private;
    output r2 as u64.private;`;

app.post('/prove', async (req, res) => {
    try {
        const { market_id, amount, outcome } = req.body;
        const account = new Account();
        const networkClient = new AleoNetworkClient("https://api.explorer.provable.com/v1");
        const recordProvider = new NetworkRecordProvider(account, networkClient);
        const programManager = new ProgramManager(
            "https://api.explorer.provable.com/v1",
            undefined,
            recordProvider
        );
        programManager.setAccount(account);

        const result = await programManager.run(
            PROGRAM,
            "place_bet",
            [
                `${market_id}u64`,
                outcome === "YES" ? "true" : "false",
                `${amount}u64`,
                account.address().to_string(),
            ],
            false
        );

        console.log("Full result:", JSON.stringify(result));
        console.log("Result type:", typeof result);
        console.log("Result keys:", result ? Object.keys(result) : "null");

        // Try different ways to get output
        let txId = "proof_generated";
        if (result && typeof result.getOutputs === "function") {
            const outputs = result.getOutputs();
            console.log("Outputs:", outputs);
            txId = outputs?.[0] ?? txId;
        } else if (result && result.outputs) {
            txId = result.outputs[0] ?? txId;
        } else if (typeof result === "string") {
            txId = result;
        }

        res.json({ tx_id: txId });
    } catch (err) {
        console.error("Prover error full:", err);
        res.status(500).json({ error: err?.message ?? String(err) });
    }
});

app.listen(3001, () => console.log('Aleo prover running on :3001'));
