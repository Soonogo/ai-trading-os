'use client';

import { Send, Bot, User } from 'lucide-react';
import { useState } from 'react';
import { Button } from '@aitos/design-system';
import { useChat } from '@/stores/chat-store';
import { cn } from '@/lib/utils';

export function AIChat() {
  const [input, setInput] = useState('');
  const { messages, isTyping, addMessage, setTyping } = useChat();

  const send = async () => {
    if (!input.trim()) return;
    const traceId = `trace-${Date.now()}`;
    const userMsg = {
      id: `${traceId}-u`,
      role: 'user' as const,
      content: input,
      traceId,
      timestamp: Date.now(),
    };
    addMessage(userMsg);
    setInput('');
    setTyping(true);

    try {
      const res = await fetch('/api/v1/prompts', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          trace_id: traceId,
          prompt: input,
          user_id: 'user-1',
        }),
      });
      const data = await res.json();
      addMessage({
        id: `${traceId}-a`,
        role: 'agent',
        content: `Dispatched prompt to agents. Trace ID: ${data.trace_id}`,
        traceId,
        timestamp: Date.now(),
      });
    } catch (err) {
      addMessage({
        id: `${traceId}-e`,
        role: 'system',
        content: `Could not reach AI services: ${(err as Error).message}`,
        traceId,
        timestamp: Date.now(),
      });
    } finally {
      setTyping(false);
    }
  };

  return (
    <div className="flex h-full flex-col border-l border-white/[0.06] glass-elevated">
      <div className="flex h-16 items-center border-b border-white/[0.06] px-4">
        <Bot className="mr-2 h-5 w-5 text-primary" />
        <span className="font-semibold">AI Copilot</span>
      </div>

      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.map((msg) => (
          <div
            key={msg.id}
            className={cn(
              'flex gap-3',
              msg.role === 'user' ? 'flex-row-reverse' : 'flex-row'
            )}
          >
            <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-secondary text-xs font-medium">
              {msg.role === 'user' ? <User className="h-4 w-4" /> : <Bot className="h-4 w-4" />}
            </div>
            <div
              className={cn(
                'max-w-[80%] rounded-xl px-3 py-2 text-sm',
                msg.role === 'user'
                  ? 'bg-primary text-primary-foreground'
                  : 'glass text-foreground'
              )}
            >
              {msg.content}
            </div>
          </div>
        ))}
        {isTyping && (
          <div className="flex gap-3">
            <div className="h-7 w-7 rounded-full bg-secondary" />
            <div className="glass rounded-xl px-3 py-2 text-sm text-muted-foreground">
              Agents are thinking…
            </div>
          </div>
        )}
      </div>

      <div className="border-t border-white/[0.06] p-3">
        <form
          onSubmit={(e) => {
            e.preventDefault();
            send();
          }}
          className="flex items-center gap-2"
        >
          <input
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Ask the agents…"
            className="flex-1 rounded-lg border border-input bg-background/60 px-3 py-2 text-sm outline-none ring-0 placeholder:text-muted-foreground focus:border-primary"
          />
          <Button type="submit" size="icon" variant="glass" aria-label="Send">
            <Send className="h-4 w-4" />
          </Button>
        </form>
      </div>
    </div>
  );
}
