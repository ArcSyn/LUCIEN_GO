from agents.agent_base import Agent

class Rewriter(Agent):
    def run(self, input_text: str) -> str:
        return f"[Rewriter] I received: '{input_text}'"

# Allow direct CLI calls
if __name__ == "__main__":
    agent = Rewriter()
    print(agent.run("Hello from CLI"))