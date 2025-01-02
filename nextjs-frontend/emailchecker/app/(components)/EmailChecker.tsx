"use client";
import React, { useState } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Check, X } from "lucide-react";

const DomainValidator = () => {
  const [domain, setDomain] = useState("");
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const response = await fetch("http://localhost:8080/form", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ domainurl: domain }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      console.log("Response data:", data); // Debug log
      setResults(data);
    } catch (err) {
      console.error("Error:", err); // Debug log
      setError(`Failed to validate domain: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto p-4">
      <Card>
        <CardHeader>
          <CardTitle>Domain Validator</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="flex gap-2">
              <Input
                type="text"
                value={domain}
                onChange={(e) => setDomain(e.target.value)}
                placeholder="Enter domain (e.g., example.com)"
                className="flex-1"
              />
              <Button type="submit" disabled={loading}>
                {loading ? "Validating..." : "Validate"}
              </Button>
            </div>

            {error && <div className="text-red-500">{error}</div>}

            <div className="space-y-4">
              {results.map((result, index) => (
                <Card key={index} className="p-4">
                  <h3 className="font-bold mb-2">{result.domain}</h3>
                  <div className="grid grid-cols-2 gap-2">
                    <div className="flex items-center gap-2">
                      {result.hasMX ? (
                        <Check className="text-green-500" />
                      ) : (
                        <X className="text-red-500" />
                      )}
                      <span>MX Records</span>
                    </div>
                    <div className="flex items-center gap-2">
                      {result.hasSPF ? (
                        <Check className="text-green-500" />
                      ) : (
                        <X className="text-red-500" />
                      )}
                      <span>SPF Record</span>
                    </div>
                    <div className="flex items-center gap-2">
                      {result.hasDMARC ? (
                        <Check className="text-green-500" />
                      ) : (
                        <X className="text-red-500" />
                      )}
                      <span>DMARC Record</span>
                    </div>
                  </div>

                  {result.spfRecord && (
                    <div className="mt-2">
                      <div className="font-semibold">SPF Record:</div>
                      <div className="text-sm bg-gray-50 p-2 rounded">
                        {result.spfRecord}
                      </div>
                    </div>
                  )}

                  {result.dmarcRecord && (
                    <div className="mt-2">
                      <div className="font-semibold">DMARC Record:</div>
                      <div className="text-sm bg-gray-50 p-2 rounded">
                        {result.dmarcRecord}
                      </div>
                    </div>
                  )}
                </Card>
              ))}
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};

export default DomainValidator;
