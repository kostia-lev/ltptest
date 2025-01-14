import { useState } from "react";
import axios from "axios";

const LTP = () => {
    const [pairs, setPairs] = useState("BTCUSD,BTCCHF,BTCEUR");
    const [data, setData] = useState<{pair: string, amount: number}[]>([]);
    const [error, setError] = useState("");

    const fetchLTP = async () => {
        try {
            setError("");
            const response = await axios.get(
                `http://localhost:8082/api/v1/ltp?pairs=${pairs}`
            );
            setData(response.data.ltp);
        } catch {
            setError("Failed to fetch LTP data.");
        }
    };

    return (
        <div>
            <h1>Bitcoin Last Traded Price</h1>
            <input
                type="text"
                value={pairs}
                onChange={(e) => setPairs(e.target.value)}
                placeholder="Enter currency pairs (comma-separated)"
            />
            <button onClick={fetchLTP}>Fetch LTP</button>

            {error && <p style={{ color: "red" }}>{error}</p>}

            <table style={{ marginTop: "20px" }}>
                <thead>
                <tr>
                    <th>Pair</th>
                    <th>Amount</th>
                </tr>
                </thead>
                <tbody>
                {data.map((item) => (
                    <tr key={item.pair}>
                        <td>{item.pair}</td>
                        <td>{item.amount}</td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
};

export default LTP;
