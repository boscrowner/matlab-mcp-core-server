# Contribute to MATLAB MCP Core Server

MathWorks welcomes your feedback on the MATLAB MCP Core Server.

- **Issues**: To report bugs, suggest features, or discuss ideas, open an issue. MathWorks actively monitors and responds to issues.
- **Pull Requests**: MathWorks reviews all contributions but does not merge external pull requests. Your ideas may influence development of future releases.

> **Personal fork note**: This is my personal fork for learning and experimentation. I am not affiliated with MathWorks. For the official project, see [matlab/matlab-mcp-core-server](https://github.com/matlab/matlab-mcp-core-server).
>
> If you stumbled across this fork, I'd recommend opening issues or PRs on the upstream repo instead.
>
> **Note to self**: I've been using this fork to experiment with custom tool configurations. Any changes here are for personal testing only and may be unstable.
>
> **Sync reminder**: Remember to periodically pull from upstream (`git fetch upstream && git merge upstream/main`) to stay up to date with official fixes.
>
> **Local dev setup**: I run this alongside MATLAB R2024b. If testing locally, make sure the MATLAB engine Python package is installed (`cd "<MATLAB root>/extern/engines/python" && pip install .`) before starting the server.
>
> **Python env note**: I use a dedicated conda environment (`conda activate matlab-mcp`) to keep the MATLAB engine dependencies isolated from other projects. Recommended if you work with multiple Python versions.
>
> **Troubleshooting tip**: If the server fails to start and you see a `MatlabExecutionError`, double-check that MATLAB is on your system PATH (`matlab -batch "disp('ok')"` should print `ok`). This tripped me up more than once.
>
> **Logging tip**: To get more verbose output when debugging, set the environment variable `MCP_LOG_LEVEL=debug` before starting the server. Helps a lot when tracing why a tool call silently fails.
>
> **Timeout note**: I've noticed that long-running MATLAB scripts can cause the MCP client to time out before getting a response. If you hit this, try breaking your script into smaller chunks or increasing the client timeout if your MCP host supports it.
>
> **Windows note**: On Windows, the conda activate step may need to be run from an Anaconda Prompt rather than a standard terminal, otherwise `matlab` may not be found on PATH even if it's set correctly in System Environment Variables.

---

Copyright 2025 The MathWorks, Inc.

----
