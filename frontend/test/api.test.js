import assert from "node:assert/strict";
import { test } from "node:test";
import { ApiError, createApiClient, requestJson } from "../src/api.js";

test("requestJson sends JSON with cookies enabled", async () => {
  const calls = [];
  const fetcher = async (url, options) => {
    calls.push({ url, options });
    return new Response("", { status: 201 });
  };

  const result = await requestJson(fetcher, "", "/api/add/", {
    method: "POST",
    body: { text: "Buy milk" },
  });

  assert.equal(result, null);
  assert.equal(calls[0].url, "/api/add/");
  assert.equal(calls[0].options.credentials, "include");
  assert.equal(calls[0].options.method, "POST");
  assert.equal(calls[0].options.headers.get("Content-Type"), "application/json");
  assert.equal(calls[0].options.body, JSON.stringify({ text: "Buy milk" }));
});

test("requestJson exposes backend error messages", async () => {
  const fetcher = async () =>
    new Response(JSON.stringify({ message: "Unauthorized user" }), { status: 401 });

  await assert.rejects(
    requestJson(fetcher, "", "/api/"),
    (error) => {
      assert.ok(error instanceof ApiError);
      assert.equal(error.status, 401);
      assert.equal(error.message, "Unauthorized user");
      return true;
    },
  );
});

test("protected requests refresh once after an unauthorized response", async () => {
  const calls = [];
  const fetcher = async (url) => {
    calls.push(url);

    if (calls.length === 1) {
      return new Response(JSON.stringify({ message: "Unauthorized user" }), { status: 401 });
    }

    if (url === "/api/refresh/") {
      return new Response("", { status: 200 });
    }

    return new Response(JSON.stringify([{ id: 1, user_id: 2, text: "Ship frontend" }]), {
      status: 200,
    });
  };

  const api = createApiClient({ fetcher });
  const notes = await api.getNotes();

  assert.deepEqual(calls, ["/api/", "/api/refresh/", "/api/"]);
  assert.deepEqual(notes, [{ id: 1, user_id: 2, text: "Ship frontend" }]);
});
